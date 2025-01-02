package gen

import (
	"fmt"
	"go/ast"

	"github.com/spf13/cobra"
	"golang.org/x/tools/go/packages"
)

func NewGenCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gen",
		Short: "Generate code",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := &packages.Config{
				Mode: packages.NeedName | packages.NeedFiles | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedSyntax,
			}

			pkgs, err := packages.Load(cfg, "./...")
			if err != nil {
				return err
			}

			for _, pkg := range pkgs {
				fmt.Println("Package:", pkg.Name)
				for _, file := range pkg.Syntax {
					fmt.Println("\tFile:", file.Name)
					ast.Inspect(file, func(n ast.Node) bool {
						switch x := n.(type) {
						case *ast.FuncDecl:
							fmt.Println("\t\tFunction Name:", x.Name.Name)
							if x.Recv != nil {
								for _, recv := range x.Recv.List {
									fmt.Println("\t\tReceiver:", recv)
									fmt.Println("\t\tReceiver Type:", recv.Type)
								}
							}
							fmt.Println("\t\t----------------")
						case *ast.GenDecl:
							fmt.Println("\t\tGenDecl:", x.Tok)
							for _, spec := range x.Specs {
								switch s := spec.(type) {
								case *ast.TypeSpec:
									fmt.Println("\t\tTypeSpec:", s.Name)
									fmt.Println("\t\tTypeSpec Type:", s.Type)
								case *ast.ValueSpec:
									fmt.Println("\t\tValueSpec:", s.Names)
									fmt.Println("\t\tValueSpec Type:", s.Type)
									fmt.Println("\t\tValueSpec Values:", s.Values)
								}
							}
							if x.Doc != nil {
								fmt.Println("\t\tGenDecl Doc:", x.Doc.Text())
							}
							fmt.Println("\t\t----------------")
						}

						return true
					})

					fmt.Println("\t----------------")
				}

				fmt.Println("\t----------------")
			}

			return nil
		},
	}

	return cmd
}
