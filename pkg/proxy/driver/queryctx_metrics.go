package driver

import (
	"context"
	gomysql "github.com/siddontang/go-mysql/mysql"

	"github.com/pingcap/parser/ast"
	"github.com/tidb-incubator/weir/pkg/proxy/metrics"
	wast "github.com/tidb-incubator/weir/pkg/util/ast"
)

func (q *QueryCtxImpl) recordQueryMetrics(ctx context.Context, stmt ast.StmtNode, sqlResult *gomysql.Result, err error, durationMilliSecond float64) {
	ns := q.ns.Name()
	db := q.currentDB
	firstTableName, _ := wast.GetAstTableNameFromCtx(ctx)
	stmtType := metrics.GetStmtTypeName(stmt)
	retLabel := metrics.RetLabel(err)

	metrics.QueryCtxQueryCounter.WithLabelValues(ns, db, firstTableName, stmtType, retLabel).Inc()
	metrics.QueryCtxQueryDurationHistogram.WithLabelValues(ns, db, firstTableName, stmtType).Observe(durationMilliSecond)
	if sqlResult != nil {
		metrics.QueryCtxQuerySizeCounter.WithLabelValues(ns, db, firstTableName, stmtType).Add(float64(len(sqlResult.RawPkg)))
	}
}

func (q *QueryCtxImpl) recordDeniedQueryMetrics(ctx context.Context, stmt ast.StmtNode) {
	ns := q.ns.Name()
	db := q.currentDB
	firstTableName, _ := wast.GetAstTableNameFromCtx(ctx)
	stmtType := metrics.GetStmtTypeName(stmt)

	metrics.QueryCtxQueryDeniedCounter.WithLabelValues(ns, db, firstTableName, stmtType).Inc()
}
