package store

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func withTestDB(t *testing.T, fn func(*sqlx.DB, sqlmock.Sqlmock)) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer mockDB.Close()

	db := sqlx.NewDb(mockDB, "sqlmock")
	fn(db, mock)
}

func TestCreateProduct(t *testing.T) {
	p := &Product{
		Name:         "test product",
		Image:        "test.jpg",
		Category:     "test category",
		Description:  "test description",
		Rating:       5,
		NumReviews:   10,
		Price:        100.0,
		CountInStock: 100,
	}

	tcs := []struct {
		name string
		test func(*testing.T, *MySQLStore, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, st *MySQLStore, mock sqlmock.Sqlmock) {
				query := `
					INSERT INTO products (name, image, category, description, rating, num_reviews, price, count_in_stock) 
					VALUES (?, ?, ?, ?, ?, ?, ?, ?)
				`
				mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(1, 1))
				cp, err := st.CreateProduct(context.Background(), p)

				require.NoError(t, err)
				require.Equal(t, int64(1), cp.ID)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed inserting product",
			test: func(t *testing.T, st *MySQLStore, mock sqlmock.Sqlmock) {
				query := `
					INSERT INTO products (name, image, category, description, rating, num_reviews, price, count_in_stock) 
					VALUES (?, ?, ?, ?, ?, ?, ?, ?)
				`
				mock.ExpectExec(query).WillReturnError(fmt.Errorf("error inserting product"))
				_, err := st.CreateProduct(context.Background(), p)

				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed getting last insert id",
			test: func(t *testing.T, st *MySQLStore, mock sqlmock.Sqlmock) {
				query := `
					INSERT INTO products (name, image, category, description, rating, num_reviews, price, count_in_stock) 
					VALUES (?, ?, ?, ?, ?, ?, ?, ?)
				`
				mock.ExpectExec(query).WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("error getting last insert id")))
				_, err := st.CreateProduct(context.Background(), p)

				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				st := NewMySQLStore(db)
				tc.test(t, st, mock)
			})
		})
	}
}

func TestGetProduct(t *testing.T) {
	p := &Product{
		Name:         "test product",
		Image:        "test.jpg",
		Category:     "test category",
		Description:  "test description",
		Rating:       5,
		NumReviews:   10,
		Price:        100.0,
		CountInStock: 100,
	}

	tcs := []struct {
		name string
		test func(*testing.T, *MySQLStore, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, st *MySQLStore, mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "image", "category", "description", "rating", "num_reviews", "price", "count_in_stock", "created_at", "updated_at"}).
					AddRow(1, p.Name, p.Image, p.Category, p.Description, p.Rating, p.NumReviews, p.Price, p.CountInStock, p.CreatedAt, p.UpdatedAt)

				mock.ExpectQuery("SELECT * FROM products WHERE id=?").WithArgs(1).WillReturnRows(rows)

				gp, err := st.GetProduct(context.Background(), 1)
				require.NoError(t, err)
				require.Equal(t, int64(1), gp.ID)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed getting product",
			test: func(t *testing.T, st *MySQLStore, mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT * FROM products WHERE id=?").WithArgs(1).WillReturnError(fmt.Errorf("error getting product"))

				_, err := st.GetProduct(context.Background(), 1)
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				st := NewMySQLStore(db)
				tc.test(t, st, mock)
			})
		})
	}
}

func TestListProducts(t *testing.T) {
	p := &Product{
		Name:         "test product",
		Image:        "test.jpg",
		Category:     "test category",
		Description:  "test description",
		Rating:       5,
		NumReviews:   10,
		Price:        100.0,
		CountInStock: 100,
	}

	tcs := []struct {
		name string
		test func(*testing.T, *MySQLStore, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, st *MySQLStore, mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "image", "category", "description", "rating", "num_reviews", "price", "count_in_stock", "created_at", "updated_at"}).
					AddRow(1, p.Name, p.Image, p.Category, p.Description, p.Rating, p.NumReviews, p.Price, p.CountInStock, p.CreatedAt, p.UpdatedAt)

				mock.ExpectQuery("SELECT * FROM products").WillReturnRows(rows)

				products, err := st.ListProducts(context.Background())
				require.NoError(t, err)
				require.Len(t, products, 1)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed listing products",
			test: func(t *testing.T, st *MySQLStore, mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT * FROM products").WillReturnError(fmt.Errorf("error listing products"))

				_, err := st.ListProducts(context.Background())
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				st := NewMySQLStore(db)
				tc.test(t, st, mock)
			})
		})
	}
}

func TestUpdateProduct(t *testing.T) {
	p := &Product{
		ID:           1,
		Name:         "test product",
		Image:        "test.jpg",
		Category:     "test category",
		Description:  "test description",
		Rating:       5,
		NumReviews:   10,
		Price:        100.0,
		CountInStock: 100,
	}

	np := &Product{
		ID:           1,
		Name:         "new test product",
		Image:        "new test.jpg",
		Category:     "new test category",
		Description:  "new test description",
		Rating:       4,
		NumReviews:   20,
		Price:        200.0,
		CountInStock: 200,
	}

	tcs := []struct {
		name string
		test func(*testing.T, *MySQLStore, sqlmock.Sqlmock)
	}{
		{
			name: "success",
			test: func(t *testing.T, st *MySQLStore, mock sqlmock.Sqlmock) {
				queryCreate := `
					INSERT INTO products (name, image, category, description, rating, num_reviews, price, count_in_stock) 
					VALUES (?, ?, ?, ?, ?, ?, ?, ?)
				`
				mock.ExpectExec(queryCreate).WillReturnResult(sqlmock.NewResult(1, 1))
				cp, err := st.CreateProduct(context.Background(), p)
				require.NoError(t, err)
				require.Equal(t, int64(1), cp.ID)

				queryUpdate := `
					UPDATE products 
					SET name=?, image=?, category=?, description=?, rating=?, num_reviews=?, price=?, count_in_stock=?, updated_at=?
					WHERE id=?
				`
				mock.ExpectExec(queryUpdate).WillReturnResult(sqlmock.NewResult(1, 1))

				up, err := st.UpdateProduct(context.Background(), np)
				require.NoError(t, err)
				require.Equal(t, int64(1), up.ID)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed updating product",
			test: func(t *testing.T, st *MySQLStore, mock sqlmock.Sqlmock) {
				query := `
					UPDATE products 
					SET name=?, image=?, category=?, description=?, rating=?, num_reviews=?, price=?, count_in_stock=?, updated_at=?
					WHERE id=?
				`
				mock.ExpectExec(query).WillReturnError(fmt.Errorf("error updating product"))

				_, err := st.UpdateProduct(context.Background(), np)
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				st := NewMySQLStore(db)
				tc.test(t, st, mock)
			})
		})
	}
}

func TestDeleteProduct(t *testing.T) {
	tcs := []struct {
		name string
		test func(*testing.T, *MySQLStore, sqlmock.Sqlmock)
	} {
		{
			name: "success",
			test: func(t *testing.T, st *MySQLStore, mock sqlmock.Sqlmock) {
				query := `DELETE FROM products WHERE id=?`
				mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(1, 1))

				err := st.DeleteProduct(context.Background(), 1)
				require.NoError(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
		{
			name: "failed deleting product",
			test: func(t *testing.T, st *MySQLStore, mock sqlmock.Sqlmock) {
				query := `DELETE FROM products WHERE id=?`
				mock.ExpectExec(query).WillReturnError(fmt.Errorf("error deleting product"))

				err := st.DeleteProduct(context.Background(), 1)
				require.Error(t, err)

				err = mock.ExpectationsWereMet()
				require.NoError(t, err)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			withTestDB(t, func(db *sqlx.DB, mock sqlmock.Sqlmock) {
				st := NewMySQLStore(db)
				tc.test(t, st, mock)
			})
		})
	}
}