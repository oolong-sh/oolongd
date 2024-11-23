# Oolong Development Notes
Oct 2, 2024

## Tokenizing -> NGram Flow

- Frequency Score (tf-idf)
    - requires number of occurrences in document
    - number of occurrences across documents
    - requires number of documents
    - TF should be calculated at the document level
    - IDF should be calculated after all documents are fully processed

- Token/NGram info:
    - Need to store all tokens in a document (and Ngrams)
    - Store filepaths for each document
    - Need occurences of ngrams in each document and in all documents

## API

- Build API for frontend-backend communication
    - CRUD endpoints for files
    - Weight recalculation trigger
    - Graph data getter
