package main

import (
	"flag"
	"fmt"
	"github.com/vektah/gqlparser/v2"
	"os"
	"sort"
	"strings"

	"github.com/vektah/gqlparser/v2/ast"
)

// Configuration holds all command-line parameters
type Config struct {
	SchemaFilePath string
	OutputPath     string
	Header         string
	ClientName     string
	ErrorType      string
}

func main() {
	// Parse command-line flags
	config := parseFlags()

	// Read schema file
	schemaFile, err := os.ReadFile(config.SchemaFilePath)
	if err != nil {
		fmt.Printf("Error reading schema file: %v\n", err)
		os.Exit(1)
	}

	// Create output file
	outputFile, err := os.Create(config.OutputPath)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer outputFile.Close()

	// Process header content
	headerContent := strings.ReplaceAll(config.Header, `\n`, "\n")

	// Generate TypeScript code
	createLibrary(schemaFile, outputFile, headerContent, config.ClientName, config.ErrorType)
}

func parseFlags() Config {
	// Define command-line flags
	schemaFilePath := flag.String("schema", "", "Path to GraphQL schema file")
	outputPath := flag.String("output", "", "Path to output TypeScript file")
	header := flag.String("header", "", "Header code to insert at the beginning of the TypeScript file")
	clientName := flag.String("client", "", "Name of the client (e.g., apolloClient or client)")
	errorType := flag.String("error", "", "Return errors as return value")

	// Parse flags
	flag.Parse()

	// Check if all required flags are set
	if *schemaFilePath == "" || *outputPath == "" || *header == "" || *clientName == "" {
		fmt.Println("All parameters are required.")
		flag.Usage()
		os.Exit(1)
	}

	return Config{
		SchemaFilePath: *schemaFilePath,
		OutputPath:     *outputPath,
		Header:         *header,
		ClientName:     *clientName,
		ErrorType:      *errorType,
	}
}

func createLibrary(schemaFile []byte, exportFile *os.File, header string, clientName string, errorType string) {
	// Parse the GraphQL schema
	schemaDocument, err := gqlparser.LoadSchema(&ast.Source{
		Input: string(schemaFile),
	})
	if err != nil {
		panic(err)
	}

	var outputTypesMap = make(map[string]string)
	var outputFuncMap = make(map[string]string)

	funcListMutation := []string{}
	funcListQuery := []string{}

	// Process schema types and fields
	for _, typeDef := range schemaDocument.Types {
		outputTypes := ""

		// Skip internal and scalar types
		if (len(typeDef.Name) >= 2 && typeDef.Name[0:2] == "__") ||
			typeDef.Name == "Boolean" || typeDef.Name == "Float" ||
			typeDef.Name == "Int" || typeDef.Name == "String" ||
			typeDef.Name == "ID" {
			continue
		}

		// Process Query and Mutation types
		if typeDef.Name == "Query" || typeDef.Name == "Mutation" {
			clientFunc := "query"
			clientFuncGql := "query"
			if typeDef.Name == "Mutation" {
				clientFunc = "mutate"
				clientFuncGql = "mutation"
			}

			for _, field := range typeDef.Fields {
				outputFunc := ""

				// Skip internal fields
				if (len(field.Name) >= 2 && field.Name[0:2] == "__") ||
					field.Name == "Boolean" || field.Name == "Float" ||
					field.Name == "Int" || field.Name == "String" ||
					field.Name == "ID" {
					continue
				}

				// Process arguments
				args, input, line1, line2 := processArguments(field.Arguments)

				// Generate query response
				queryResponse := toGraphqlResponse(schemaDocument, field.Type.Name(), 5)

				// Track function names by type
				if typeDef.Name == "Query" {
					funcListQuery = append(funcListQuery, field.Name)
				} else {
					funcListMutation = append(funcListMutation, field.Name)
				}

				// Generate function code based on error handling preference
				outputFunc = generateFunctionCode(field, args, clientFunc, clientFuncGql,
					input, line2, line1, clientName,
					queryResponse, errorType)

				outputFuncMap[field.Name] = outputFunc
			}
			continue
		}

		// Skip Mutation type (already handled above)
		if typeDef.Name == "Mutation" {
			continue
		}

		// Generate interface for regular types
		outputTypes += "export interface " + typeDef.Name + " {\n"
		for _, field := range typeDef.Fields {
			outputTypes += "\t" + field.Name + ": " + getType(field.Type.String()) + ";\n"
		}
		outputTypes += "}\n\n"

		outputTypesMap[typeDef.Name] = outputTypes
	}

	// Write to output file
	writeToOutputFile(exportFile, header, outputTypesMap, outputFuncMap, funcListQuery, funcListMutation)
}

func processArguments(arguments ast.ArgumentDefinitionList) (string, string, string, string) {
	args := ""
	input := ""
	line1 := ""
	line2 := ""

	for i, arg := range arguments {
		comma := ", "
		breakLine := ",\n"
		if i == len(arguments)-1 {
			comma = ""
			breakLine = ""
		}
		args += fmt.Sprintf("%s: %s%s", arg.Name, getType(arg.Type.String()), comma)
		input += fmt.Sprintf("\t\t\t\t%s: %s%s", arg.Name, arg.Name, breakLine)
		line2 += fmt.Sprintf("%s: $%s%s", arg.Name, arg.Name, comma)
		line1 += fmt.Sprintf("$%s: %s%s", arg.Name, arg.Type.String(), comma)
	}

	if line2 != "" {
		line2 = "(" + line2 + ")"
	}

	if line1 != "" {
		line1 = "(" + line1 + ")"
	}

	return args, input, line1, line2
}

func generateFunctionCode(field *ast.FieldDefinition, args, clientFunc, clientFuncGql,
	input, line2, line1, clientName, queryResponse, errorType string) string {

	outputFunc := ""
	tmp := ""

	if errorType == "" {
		outputFunc += fmt.Sprintf("async function %s(%s): Promise<%s> {",
			field.Name, args, getType(field.Type.String()))
		tmp = fmt.Sprintf(`
const response = await %[8]s.%[3]s({
					%[4]s: gql('
				 %[4]s %[1]s%[7]s {
					%[1]s%[6]s %[2]s
				}
				'),
					variables: {
%[5]s
					}
				});

				// Überprüfe, ob die Antwort erfolgreich war
				return response.data.%[1]s;`,
			field.Name, queryResponse, clientFunc, clientFuncGql, input, line2, line1, clientName)
	} else {
		outputFunc += fmt.Sprintf("async function %s(%s): Promise<[%s | null, any | null]> {",
			field.Name, args, getType(field.Type.String()))
		tmp = fmt.Sprintf(`
	try {
		const response = await %[8]s.%[3]s({
			%[4]s: gql('
				 %[4]s %[1]s%[7]s {
					%[1]s%[6]s %[2]s
				}
				'),
			variables: {
				%[5]s
			}
		});
	
		if (response.errors && response.errors.length > 0) {
			return [null, response.errors];
		}
	
		if (!response.data || !response.data.%[1]s) {
			return [null, new Error('Invalid response structure')];
		}
	
		return [response.data.%[1]s, null];
	} catch (error) {
		return [null, error];
	}`,
			field.Name, queryResponse, clientFunc, clientFuncGql, input, line2, line1, clientName)
	}

	tmp = strings.ReplaceAll(tmp, "'", "`")
	outputFunc += tmp
	outputFunc += "\n}\n\n"

	return outputFunc
}

func writeToOutputFile(exportFile *os.File, header string, outputTypesMap, outputFuncMap map[string]string,
	funcListQuery, funcListMutation []string) {

	// Write header
	exportFile.WriteString(header + "\n\n")

	// Write types in sorted order
	writeMapInSortedOrder(exportFile, outputTypesMap)

	// Write functions in sorted order
	writeMapInSortedOrder(exportFile, outputFuncMap)

	// Write query function list
	writeFunctionList(exportFile, "query", funcListQuery)

	// Write mutation function list
	writeFunctionList(exportFile, "mutation", funcListMutation)

	// Add final newlines
	exportFile.WriteString("\n\n")
}

func writeMapInSortedOrder(file *os.File, dataMap map[string]string) {
	var keys []string
	for k := range dataMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		file.WriteString(dataMap[k])
	}
}

func writeFunctionList(file *os.File, listName string, functions []string) {
	file.WriteString(fmt.Sprintf("export const %s = {\n", listName))

	sort.Strings(functions)
	for i, f := range functions {
		comma := ",\n"
		if i == len(functions)-1 {
			comma = ""
		}
		file.WriteString(fmt.Sprintf("\t%s%s", f, comma))
	}

	file.WriteString("\n};\n")
}

func getType(typ string) string {
	// Replace all ! with empty string
	typ = strings.ReplaceAll(typ, "!", "")

	// Convert GraphQL types to TypeScript types
	switch typ {
	case "ID", "String":
		return "string"
	case "Int":
		return "number"
	case "Float":
		return "number"
	case "Boolean":
		return "boolean"
	default:
		if typ[0] == '[' && typ[len(typ)-1] == ']' {
			return fmt.Sprintf("%s[]", getType(typ[1:len(typ)-1]))
		} else {
			return typ
		}
	}
}

func indentTabs(n int) string {
	ret := ""
	for i := 0; i < n; i++ {
		ret += "\t"
	}
	return ret
}

func toGraphqlResponse(schemaDocument *ast.Schema, responseTypeName string, indent int) string {
	if ok := schemaDocument.Types[responseTypeName].Fields; ok != nil {
		response := "{\n"
		for _, field := range schemaDocument.Types[responseTypeName].Fields {
			response += indentTabs(indent+1) + field.Name + " "
			if ok := schemaDocument.Types[field.Type.Name()].Fields; ok != nil {
				response += toGraphqlResponse(schemaDocument, field.Type.Name(), indent+1)
			}
			response += "\n"
		}
		response += indentTabs(indent) + "}"
		return response
	}
	return ""
}
