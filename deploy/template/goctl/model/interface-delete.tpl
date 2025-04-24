Delete(ctx context.Context, {{.lowerStartCamelPrimaryKey}} {{.dataType}}) error
DeleteWithSession(ctx context.Context, session sqlx.Session, {{.lowerStartCamelPrimaryKey}} {{.dataType}}) error