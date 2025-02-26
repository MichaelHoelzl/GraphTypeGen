go run ../main.go  -schema="./example.graphql" \
                 -output="./output.ts" \
                 -header="import { gql } from '@apollo/client';" \
                 -client="apolloClient" \
                 -error="true"