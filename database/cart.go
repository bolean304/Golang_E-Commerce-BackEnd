package database

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/bolean304/e-commerce-cart/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrCantFindProduct    = errors.New("can't find the product")
	ErrCantDecodeProduct  = errors.New("can't decode the product")
	ErrUserIdIsNotValid   = errors.New("this user is not valid")
	ErrCantUpdateUser     = errors.New("can't add this product to the cart")
	ErrCantRemoveItemCart = errors.New("can't remove this item from the cart")
	ErrCantGetItem        = errors.New("was unable to get the item from the cart")
	ErrCantBuyCartItem    = errors.New("can't update the purchase of the item")
)

func AddProductToCart(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	searchFromDB := prodCollection.FindOne(ctx, bson.M{"_id": productID})
	if searchFromDB.Err() != nil {
		log.Println(searchFromDB.Err())
		return ErrCantFindProduct
	}
	var productCart []models.ProductUser
	err := searchFromDB.Decode(&productCart)
	if err != nil {
		log.Println(err)
		return ErrCantDecodeProduct
	}

	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{
		{Key: "$push", Value: bson.D{
			{Key: "usercart", Value: bson.M{"$each": productCart}},
		}},
	}

	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return ErrCantUpdateUser
	}
	return nil
}

func RemoveCartItem(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.M{"$pull": bson.M{
		"usercart": bson.M{
			"_id": productID,
		},
	}}

	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return ErrCantRemoveItemCart
	}
	return nil
}

func BuyItemFromCart(ctx context.Context, userCollection *mongo.Collection, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	var getUserCart models.User
	var orderCart models.Order
	orderCart.Order_ID = primitive.NewObjectID()
	orderCart.Orderered_At = time.Now()
	orderCart.Order_Cart = make([]models.ProductUser, 0)
	orderCart.Payment_Method.COD = true

	unwind := bson.D{
		{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$usercart"},
		}},
	}

	grouping := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$_id"},
			{Key: "total", Value: bson.D{
				{Key: "$sum", Value: "$usercart.price"},
			}},
		}},
	}

	currentResults, err := userCollection.Aggregate(ctx, mongo.Pipeline{unwind, grouping})
	if err != nil {
		return err
	}
	defer currentResults.Close(ctx)

	var getUserCartAggregate []bson.M
	if err = currentResults.All(ctx, &getUserCartAggregate); err != nil {
		return err
	}

	var totalPrice int32
	for _, userItem := range getUserCartAggregate {
		totalPrice = userItem["total"].(int32)
	}
	orderCart.Price = int(totalPrice)

	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{
		{Key: "$push", Value: bson.D{
			{Key: "orders", Value: orderCart},
		}},
	}

	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println(err)
		return err
	}

	err = userCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: id}}).Decode(&getUserCart)
	if err != nil {
		log.Println(err)
		return err
	}

	filter2 := bson.D{{Key: "_id", Value: id}}
	update2 := bson.M{"$push": bson.M{
		"orders.$[].order_list": bson.M{"$each": getUserCart.UserCart},
	}}

	_, err = userCollection.UpdateOne(ctx, filter2, update2)
	if err != nil {
		log.Println(err)
		return err
	}

	userCartEmpty := make([]models.ProductUser, 0)
	filter3 := bson.D{primitive.E{Key: "_id", Value: id}}
	update3 := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "usercart", Value: userCartEmpty}}}}

	_, err = userCollection.UpdateOne(ctx, filter3, update3)
	if err != nil {
		return ErrCantBuyCartItem
	}
	return nil
}

func InstantBuyer(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrUserIdIsNotValid
	}

	var productDetails models.ProductUser
	var orderDetails models.Order
	orderDetails.Order_ID = primitive.NewObjectID()
	orderDetails.Orderered_At = time.Now()
	orderDetails.Order_Cart = make([]models.ProductUser, 0)
	orderDetails.Payment_Method.COD = true

	err = prodCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: productID}}).Decode(&productDetails)
	if err != nil {
		log.Println(err)
		return ErrCantFindProduct
	}

	orderDetails.Price = productDetails.Price

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{
		{Key: "$push", Value: bson.D{
			{Key: "orders", Value: orderDetails},
		}},
	}

	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println(err)
		return err
	}

	filter2 := bson.D{primitive.E{Key: "_id", Value: id}}
	update2 := bson.M{"$push": bson.M{
		"orders.$[].order_list": productDetails,
	}}

	_, err = userCollection.UpdateOne(ctx, filter2, update2)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
