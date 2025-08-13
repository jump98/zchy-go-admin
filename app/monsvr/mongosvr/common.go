package mongosvr

import (
	"context"
	"fmt"
)

// 插入数据示例
func insertDocumentData(uri, dbName, collectionName string, document interface{}) error {
	// 确保连接有效
	if err := ensureConnection(uri); err != nil {
		return fmt.Errorf("failed to ensure connection: %v", err)
	}

	// 插入数据
	collection := client.Database(dbName).Collection(collectionName)
	_, err := collection.InsertOne(context.TODO(), document)
	if err != nil {
		return fmt.Errorf("failed to insert document: %v", err)
	}

	//fmt.Println("Document inserted successfully!")
	return nil
}

// 插入数据示例
func insertArrayDocumentData(uri, dbName, collectionName string, document []interface{}) error {
	// 确保连接有效
	if err := ensureConnection(uri); err != nil {
		return fmt.Errorf("failed to ensure connection: %v", err)
	}

	// 插入数据
	collection := client.Database(dbName).Collection(collectionName)
	_, err := collection.InsertMany(context.TODO(), document)
	if err != nil {
		return fmt.Errorf("failed to insert document: %v", err)
	}

	//fmt.Println("Document inserted successfully!")
	return nil
}
