curl -X POST http://localhost:8000/api/v1/receipt/upload \
    -F "file=@/Users/sharithgodamanna/Desktop/Code/recept-app/civet/data/large-receipt-image-dataset-SRD/1001-receipt.jpg" \
    -H "Content-Type: multipart/form-data"
