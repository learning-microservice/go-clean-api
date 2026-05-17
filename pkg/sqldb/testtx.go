package sqldb

import (
	"context"
	"fmt"
)

// RunInRollbackTx はテスト専用。fn 実行後、成功・失敗にかかわらず常にロールバックする。
func (c *Client) RunInRollbackTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	return fn(withTx(ctx, tx))
}
