package sqldb

import (
	"context"
	"fmt"
)

type txFunc func(ctx context.Context) (err error)

func (c *Client) RequiredTx(ctx context.Context, txFunc txFunc) (err error) {
	// 既にトランザクションに参加している場合はオーナーシップを持たない
	if _, ok := txFromContext(ctx); ok {
		return txFunc(ctx)
	}

	// 新たにトランザクションを開始
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// トランザクションの終了処理
	// TODO: ロールバック時にエラーをログに出力
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// トランザクションをコンテキストに保存
	ctx = withTx(ctx, tx)

	// トランザクション内で処理を実行
	return txFunc(ctx)
}

func (c *Client) RequiredNewTx(ctx context.Context, txFunc txFunc) (err error) {
	// 常に新たにトランザクションを開始
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// トランザクションの終了処理
	// TODO: ロールバック時にエラーをログに出力
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
		if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// トランザクションをコンテキストに保存
	ctx = withTx(ctx, tx)

	// トランザクション内で処理を実行
	return txFunc(ctx)
}
