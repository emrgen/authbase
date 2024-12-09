#Meilisearch apis

# get all indexes

```sh
curl \
  -X GET 'http://localhost:7700/indexes' \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer 123456'
```

# create index

```sh
curl \
  -X POST 'http://localhost:7700/indexes' \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer 123456' \
  --data-binary '{
    "uid": "products",
    "primaryKey": "id"
  }'
```

# add documents 100

```sh
curl \
  -X POST 'http://localhost:7700/indexes/products/documents' \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer 123456' \
    --data-binary '[
        {"id": 1, "name": "orange", "tags": ["switch", "fruit"]},
        {"id": 2, "name": "apple", "tags": ["switch", "fruit"]},
        {"id": 3, "name": "banana", "tags": ["kill", "dog"]},
        {"id": 4, "name": "grape"},
        {"id": 5, "name": "kiwi"},
        {"id": 6, "name": "mango"},
        {"id": 7, "name": "pear"},
        {"id": 8, "name": "strawberry"},
        {"id": 9, "name": "watermelon"},
        {"id": 10, "name": "pineapple"},
        {"id": 11, "name": "blueberry"},
        {"id": 12, "name": "raspberry"},
        {"id": 13, "name": "blackberry"},
        {"id": 14, "name": "cherry"},
        {"id": 15, "name": "peach"},
        {"id": 16, "name": "plum"},
        {"id": 17, "name": "apricot"},
        {"id": 18, "name": "nectarine"},
        {"id": 19, "name": "pomegranate"},
        {"id": 20, "name": "lemon"},
        {"id": 21, "name": "lime"},
        {"id": 22, "name": "coconut"},
        {"id": 23, "name": "fig"},
        {"id": 24, "name": "date"},
        {"id": 25, "name": "persimmon"},
        {"id": 26, "name": "guava"},
        {"id": 27, "name": "papaya"},
        {"id": 28, "name": "passionfruit"},
        {"id": 29, "name": "dragonfruit"},
        {"id": 30, "name": "lychee"},
        {"id": 31, "name": "starfruit"},
        {"id": 32, "name": "kiwano"},
        {"id": 33, "name": "cantaloupe"},
        {"id": 34, "name": "honeydew"},
        {"id": 35, "name": "casaba"},
        {"id": 36, "name": "cucumber"},
        {"id": 37, "name": "zucchini"},
        {"id": 38, "name": "squash"},
        {"id": 39, "name": "pumpkin"},
        {"id": 40, "name": "tomato"},
        {"id": 41, "name": "bell pepper"},
        {"id": 42, "name": "chili pepper"},
        {"id": 43, "name": "eggplant"},
        {"id": 44, "name": "cabbage"},
        {"id": 45, "name": "broccoli"},
        {"id": 46, "name": "cauliflower"},
        {"id": 47, "name": "brussels sprout"},
        {"id": 48, "name": "kale"},
        {"id": 49, "name": "lettuce"},
        {"id": 50, "name": "spinach"},
        {"id": 51, "name": "arugula"},
        {"id": 52, "name": "watercress"},
        {"id": 53, "name": "endive"},
        {"id": 54, "name": "radicchio"},
        {"id": 55, "name": "chicory"},
        {"id": 56, "name": "dandelion"},
        {"id": 57, "name": "parsley"},
        {"id": 58, "name": "cilantro"},
        {"id": 59, "name": "basil"},
        {"id": 60, "name": "mint"},
        {"id": 61, "name": "rosemary"},
        {"id": 62, "name": "thyme"},
        {"id": 63, "name": "oregano"},
        {"id": 64, "name": "sage"},
        {"id": 65, "name": "lavender"},
        {"id": 66, "name": "lemon balm"},
        {"id": 67, "name": "chamomile"},
        {"id": 68, "name": "echinacea"},
        {"id": 69, "name": "ginger"},
        {"id": 70, "name": "turmeric"},
        {"id": 71, "name": "garlic"},
        {"id": 72, "name": "onion"},
        {"id": 73, "name": "shallot"},
        {"id": 74, "name": "scallion"},
        {"id": 75, "name": "leek"},
        {"id": 76, "name": "carrot"},
        {"id": 77, "name": "beet"},
        {"id": 78, "name": "radish"},
        {"id": 79, "name": "turnip"},
        {"id": 80, "name": "rutabaga"},
        {"id": 81, "name": "parsnip"},
        {"id": 82, "name": "potato"},
        {"id": 83, "name": "sweet potato"},
        {"id": 84, "name": "yam"},
        {"id": 85, "name": "cassava"},
        {"id": 86, "name": "taro"},
        {"id": 87, "name": "jicama"},
        {"id": 88, "name": "daikon"},
        {"id": 89, "name": "horseradish"},
        {"id": 90, "name": "wasabi"},
        {"id": 91, "name": "salsify"},
        {"id": 92, "name": "arrowroot"},
        {"id": 93, "name": "water chestnut"},
        {"id": 94, "name": "lotus root"},
        {"id": 95, "name": "burdock"},
        {"id": 96, "name": "sunchoke"},
        {"id": 97, "name": "jerusalem artichoke"},
        {"id": 98, "name": "chayote"},
        {"id": 99, "name": "okra"},
        {"id": 100, "name": "asparagus"}
    ]'
```

# add one document

```sh
curl \
  -X POST 'http://localhost:7700/indexes/products/documents' \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer 123456' \
    --data-binary '{
        "id": 101,
        "name": "cucumber",
        "tags": ["vegetable"]
    }'