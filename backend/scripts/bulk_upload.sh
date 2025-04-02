#!/bin/bash

API_URL="http://localhost:8000/api/v1/receipt/upload"
FOLDER="/Users/sharithgodamanna/Desktop/Code/recept-app/civet/backend/data/large-receipt-image-dataset-SRD"

for file in "$FOLDER"/*.jpg; do
    echo "Uploading $file..."
    curl -X POST "$API_URL" \
        -F "file=@$file" \
        -H "Content-Type: multipart/form-data"
    echo -e "\n---\n"
done
