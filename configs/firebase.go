package configs

import (
	"context"
	"fmt"
	"path/filepath"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

func SetupFirebase() (*firebase.App, context.Context, *messaging.Client) {

	ctx := context.Background()

	serviceAccountKeyFilePath, err := filepath.Abs("/Users/voidntpx/Documents/GitHub/ultima-backend/configs/serviceAccountKe.json")
	if err != nil {
		panic("Unable to load serviceAccountKeys.json file")
	}
	// fmt.Println(serviceAccountKeyFilePath)
	opt := option.WithCredentialsFile(serviceAccountKeyFilePath)

	//Firebase admin SDK initialization
	apps, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		panic("Firebase load error")
	}

	//Messaging client
	client, _ := apps.Messaging(ctx)

	fmt.Println(client)
	return apps, ctx, client
}
