// MongoDB initialization script
db = db.getSiblingDB('Agricultural');

// Create collections
db.createCollection('users');
db.createCollection('products');
db.createCollection('sales');
db.createCollection('categories');

// Create indexes
db.users.createIndex({ "email": 1 }, { unique: true });
db.products.createIndex({ "name": "text", "description": "text" });
db.products.createIndex({ "category": 1 });
db.products.createIndex({ "brand": 1 });
db.products.createIndex({ "price": 1 });
db.products.createIndex({ "stock": 1 });
db.sales.createIndex({ "product_id": 1 });
db.sales.createIndex({ "date_sold": 1 });
db.sales.createIndex({ "product_id": 1, "date_sold": 1 });
db.categories.createIndex({ "name": 1 }, { unique: true });

print('Database Agricultural initialized successfully!');
