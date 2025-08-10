package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"git.dolansoft.org/timon/sqlc-minimal-repro-playground/db"
	"github.com/jackc/pgx/v5"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://user:password@localhost:5433/testdb?sslmode=disable"
		fmt.Println("No DATABASE_URL provided, using default:", databaseURL)
		fmt.Println("Set DATABASE_URL environment variable to use a different database")
	}

	ctx := context.Background()

	conn, err := pgx.Connect(ctx, databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer conn.Close(ctx)

	if err := conn.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("✅ Successfully connected to database")

	queries := db.New(conn)

	if err := setupTestData(ctx, queries); err != nil {
		log.Fatalf("Failed to setup test data: %v", err)
	}

	if err := testFunctions(ctx, queries); err != nil {
		log.Fatalf("Test failed: %v", err)
	}

	fmt.Println("🎉 All tests passed!")
}

func setupTestData(ctx context.Context, queries *db.Queries) error {
	fmt.Println("🔧 Setting up test data...")

	ptrSandwich := db.FoodTypeSandwhich
	ptrSalad := db.FoodTypeSalad
	ptrSoup := db.FoodTypeSoup

	testFoods := []struct {
		id       int64
		food     string
		foodType *db.FoodType
	}{
		{1, "Pizza", &ptrSandwich},
		{2, "Sushi", &ptrSalad},
		{3, "Lasagne", &ptrSoup},
		{4, "Burrito", &ptrSandwich},
		{5, "Ice Cream Sundae", &ptrSalad},
		{6, "Tungsten", nil},
	}

	for _, food := range testFoods {
		err := queries.AddFood(ctx, db.AddFoodParams{
			ID:       food.id,
			Food:     food.food,
			FoodType: food.foodType,
		})
		if err != nil {
			return fmt.Errorf("failed to add food %s: %w", food.food, err)
		}
		if food.foodType == nil {
			fmt.Printf("  ✅ Added: %s (no type)\n", food.food)
		} else {
			fmt.Printf("  ✅ Added: %s (%s)\n", food.food, *food.foodType)
		}
	}

	return nil
}

func testFunctions(ctx context.Context, queries *db.Queries) error {
	fmt.Println("🧪 Testing query functions...")

	// Test 0: Search for all null food types
	fmt.Println("\n📋 Test 0: Search for all null food types (using SearchFood with nil)")
	nullFoods, err := queries.SearchFood(ctx, db.SearchFoodParams{
		Food:     "%",
		FoodType: nil,
	})
	if err != nil {
		return fmt.Errorf("failed to search for null food types: %w", err)
	}

	fmt.Printf("Found %d foods with null type:\n", len(nullFoods))
	for _, food := range nullFoods {
		foodType := "nil"
		if food.FoodType != nil {
			foodType = string(*food.FoodType)
		}
		fmt.Printf("  - ID: %d, Name: %s, Type: %s\n", food.ID, food.Food, foodType)
	}

	if len(nullFoods) != 0 {
		return fmt.Errorf("expected 0 food with null type, got %d", len(nullFoods))
	}

	// Test 1: Search for all null food types with IS NULL
	fmt.Println("\n📋 Test 1: Search for all null food types (using ListNullFoodType)")
	nullFoodsISNull, err := queries.ListNullFoodType(ctx)
	if err != nil {
		return fmt.Errorf("failed to list null food types: %w", err)
	}

	fmt.Printf("Found %d foods with null type:\n", len(nullFoodsISNull))
	for _, food := range nullFoodsISNull {
		foodType := "nil"
		if food.FoodType != nil {
			foodType = string(*food.FoodType)
		}
		fmt.Printf("  - ID: %d, Name: %s, Type: %s\n", food.ID, food.Food, foodType)
	}

	if len(nullFoodsISNull) != 1 {
		return fmt.Errorf("expected 1 food with null type, got %d", len(nullFoodsISNull))
	}

	// Test 2: Search for all sandwiches
	fmt.Println("\n📋 Test 2: Search for all sandwiches")
	sandwichType := db.FoodTypeSandwhich
	sandwiches, err := queries.SearchFood(ctx, db.SearchFoodParams{
		Food:     "%",
		FoodType: &sandwichType,
	})
	if err != nil {
		return fmt.Errorf("failed to search for sandwiches: %w", err)
	}

	fmt.Printf("Found %d sandwiches:\n", len(sandwiches))
	for _, food := range sandwiches {
		foodType := "nil"
		if food.FoodType != nil {
			foodType = string(*food.FoodType)
		}
		fmt.Printf("  - ID: %d, Name: %s, Type: %s\n", food.ID, food.Food, foodType)
	}

	if len(sandwiches) != 2 {
		return fmt.Errorf("expected 2 sandwiches, got %d", len(sandwiches))
	}

	// Test 3: Search for all salad type foods
	fmt.Println("\n📋 Test 3: Search for all salad type foods")
	saladType := db.FoodTypeSalad
	salads, err := queries.SearchFood(ctx, db.SearchFoodParams{
		Food:     "%",
		FoodType: &saladType,
	})
	if err != nil {
		return fmt.Errorf("failed to search for salads: %w", err)
	}

	fmt.Printf("Found %d salads:\n", len(salads))
	for _, food := range salads {
		foodType := "nil"
		if food.FoodType != nil {
			foodType = string(*food.FoodType)
		}
		fmt.Printf("  - ID: %d, Name: %s, Type: %s\n", food.ID, food.Food, foodType)
	}

	if len(salads) != 2 {
		return fmt.Errorf("expected 2 salads, got %d", len(salads))
	}

	// Test 4: Search for all soups containing the letter 'a'
	fmt.Println("\n📋 Test 4: Search for all soups containing the letter 'a'")
	soupType := db.FoodTypeSoup
	soupsWithA, err := queries.SearchFood(ctx, db.SearchFoodParams{
		Food:     "%a%",
		FoodType: &soupType,
	})
	if err != nil {
		return fmt.Errorf("failed to search for soups with 'a': %w", err)
	}

	fmt.Printf("Found %d soups containing 'a':\n", len(soupsWithA))
	for _, food := range soupsWithA {
		foodType := "nil"
		if food.FoodType != nil {
			foodType = string(*food.FoodType)
		}
		fmt.Printf("  - ID: %d, Name: %s, Type: %s\n", food.ID, food.Food, foodType)
	}

	if len(soupsWithA) != 1 {
		return fmt.Errorf("expected 1 soup containing 'a', got %d", len(soupsWithA))
	}

	// Test 5: Add a new food item
	fmt.Println("\n📋 Test 5: Add a new food item")
	newFoodType := db.FoodTypeSoup
	err = queries.AddFood(ctx, db.AddFoodParams{
		ID:       7,
		Food:     "Tomato Soup",
		FoodType: &newFoodType,
	})
	if err != nil {
		return fmt.Errorf("failed to add new food item: %w", err)
	}
	fmt.Println("  ✅ Added: Tomato Soup (soup)")

	// Verify the new item was added by searching for all soups
	allSoups, err := queries.SearchFood(ctx, db.SearchFoodParams{
		Food:     "%",
		FoodType: &soupType,
	})
	if err != nil {
		return fmt.Errorf("failed to verify new soup was added: %w", err)
	}

	fmt.Printf("Total soups after addition: %d\n", len(allSoups))
	if len(allSoups) != 2 {
		return fmt.Errorf("expected 2 soups after addition, got %d", len(allSoups))
	}

	return nil
}
