Draw simple overviews of cobra command trees.

Example:
```go
func generateTree(writer io.Writer, cmd *cobra.Command) {
	if err := cobratree.WriteCommandTree(writer, cmd); err != nil {
		// handle error
	}
}
```

Example default output:
```
parent
├─child
│   ├─foo
│   ├─foo1
│   │   └─bar
│   │      └─bif
│   │         └─bif2
│   └─foo2
├─child1
│   └─baz
└─child2

```