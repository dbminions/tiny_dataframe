package logicalplan

import (
	"fmt"
	"github.com/apache/arrow/go/v12/arrow"
	"strconv"
	containers "tiny_dataframe/pkg/g_containers"
)

type Expr interface {
	// DataType returns the data type of the expression. It returns error as well.
	DataType(schema containers.ISchema) (arrow.DataType, error)

	// ColumnsUsed returns the columns used in the expression.
	//TODO: replace it with ColumnsUsedExprs() []Expr
	ColumnsUsed(input LogicalPlan) []arrow.Field
	String() string
}

var _ Expr = ColumnExpr{}
var _ Expr = AliasExpr{}

var _ Expr = BooleanBinaryExpr{}
var _ Expr = MathExpr{}
var _ Expr = AggregateExpr{}

var _ Expr = LiteralStringExpr{}
var _ Expr = LiteralInt64Expr{}
var _ Expr = LiteralFloat64Expr{}

// ---------- ColumnExpr ----------

type ColumnExpr struct {
	// TODO: should this have arrow.Field?
	Name string
}

func (col ColumnExpr) DataType(schema containers.ISchema) (arrow.DataType, error) {
	for _, f := range schema.Fields() {
		if f.Name == col.Name {
			return f.Type, nil
		}
	}
	return nil, fmt.Errorf("column %s not found", col.Name)
}

func (col ColumnExpr) ColumnsUsed(input LogicalPlan) []arrow.Field {
	schema := input.Schema()
	for _, f := range schema.Fields() {
		if f.Name == col.Name {
			return []arrow.Field{f}
		}
	}
	panic(fmt.Sprintf("column %s not found", col.Name))
	return []arrow.Field{}
}

func (col ColumnExpr) String() string {
	return "#" + col.Name
}

// ---------- AliasExpr ----------

type AliasExpr struct {
	Expr  Expr
	Alias string
}

func (expr AliasExpr) DataType(schema containers.ISchema) (arrow.DataType, error) {
	return expr.Expr.DataType(schema)
}

func (expr AliasExpr) ColumnsUsed(input LogicalPlan) []arrow.Field {
	return expr.Expr.ColumnsUsed(input)
}

func (expr AliasExpr) String() string {
	return fmt.Sprintf("%s as %s", expr.Expr.String(), expr.Alias)
}

// ---------- Literals ----------

type LiteralStringExpr struct {
	Val string
}

func (lit LiteralStringExpr) DataType(schema containers.ISchema) (arrow.DataType, error) {
	return arrow.BinaryTypes.String, nil
}

func (lit LiteralStringExpr) ColumnsUsed(input LogicalPlan) []arrow.Field {
	return nil
}

func (lit LiteralStringExpr) String() string {
	return fmt.Sprintf("'%s'", lit.Val)
}

type LiteralInt64Expr struct {
	Val int64
}

func (lit LiteralInt64Expr) DataType(schema containers.ISchema) (arrow.DataType, error) {
	return arrow.PrimitiveTypes.Int64, nil
}

func (lit LiteralInt64Expr) ColumnsUsed(input LogicalPlan) []arrow.Field {
	return nil
}

func (lit LiteralInt64Expr) String() string {
	return strconv.Itoa(int(lit.Val))
}

type LiteralFloat64Expr struct {
	Val float64
}

func (lit LiteralFloat64Expr) DataType(schema containers.ISchema) (arrow.DataType, error) {
	return arrow.PrimitiveTypes.Float64, nil
}

func (lit LiteralFloat64Expr) ColumnsUsed(input LogicalPlan) []arrow.Field {
	return nil
}

func (lit LiteralFloat64Expr) String() string {
	return strconv.FormatFloat(lit.Val, 'f', -1, 64)
}
