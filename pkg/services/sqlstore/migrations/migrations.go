package migrations

import . "github.com/emretiryaki/merkut/pkg/services/sqlstore/migrator"


func AddMigrations(mg *Migrator) {
	addMigrationLogMigrations(mg)
	addAlarmMigrations(mg)
	addConditionMigrations(mg)
	addActionMigrations(mg)
	addQueryMigrations(mg)
}

func addMigrationLogMigrations(mg *Migrator) {
	migrationLog := Table{
		Name: "migration_log",
		Columns: []*Column{
			{Name: "id", Type: DB_BigInt, IsPrimaryKey: true, IsAutoIncrement: true},
			{Name: "migration_id", Type: DB_NVarchar, Length: 255},
			{Name: "sql", Type: DB_Text},
			{Name: "success", Type: DB_Bool},
			{Name: "error", Type: DB_Text},
			{Name: "timestamp", Type: DB_DateTime},
		},
	}
	mg.AddMigration("create migration_log table", NewAddTableMigration(migrationLog))
}


func addAlarmMigrations(mg *Migrator) {

	alertV1 := Table{
		Name: "alarms",
		Columns: []*Column{
			{Name: "id", Type: DB_BigInt, IsPrimaryKey: true, IsAutoIncrement: true},
			{Name: "name", Type: DB_NVarchar, Length: 250, Nullable: false},
			{Name: "state", Type: DB_NVarchar, Length: 50, Nullable: false},
			{Name: "comment", Type: DB_NVarchar,Length: 400, Nullable: false},
			{Name: "last_fired", Type: DB_DateTime, Nullable: true},
			{Name: "last_triggered", Type: DB_DateTime, Nullable: true},
			{Name: "schedule", Type: DB_NVarchar, Length: 50, Nullable: false},
			{Name: "when", Type: DB_NVarchar, Length: 100, Nullable: false},
			{Name: "indice", Type: DB_NVarchar, Length: 1000, Nullable: false},
		},
		Indices: []*Index{
			{Cols: []string{"name"}, Type: UniqueIndex},
		},
	}
	mg.AddMigration("create alarms table v1", NewAddTableMigration(alertV1))

}


func addQueryMigrations(mg *Migrator) {

	queryV1 := Table{
		Name: "queries",
		Columns: []*Column{
			{Name: "id", Type: DB_BigInt, IsPrimaryKey: true, IsAutoIncrement: true},
			{Name: "alarm_id", Type: DB_BigInt, Nullable: false},
			{Name: "field", Type: DB_NVarchar, Length: 100, Nullable: false},
			{Name: "value", Type: DB_NVarchar, Length: 100, Nullable: false},
		},
		Indices: []*Index{
			{Cols: []string{"alarm_id"}, Type: IndexType},
		},
	}
	mg.AddMigration("create queries table v1", NewAddTableMigration(queryV1))

}

func addConditionMigrations(mg *Migrator) {

	conditionV1 := Table{
		Name: "conditions",
		Columns: []*Column{
			{Name: "id", Type: DB_BigInt, IsPrimaryKey: true, IsAutoIncrement: true},
			{Name: "alarm_id", Type: DB_BigInt, Nullable: false},
			{Name: "type", Type: DB_NVarchar, Length: 50, Nullable: false},
			{Name: "operator", Type: DB_NVarchar, Length: 50, Nullable: false},
		},
		Indices: []*Index{
			{Cols: []string{"alarm_id"}, Type: IndexType},
		},
	}
	mg.AddMigration("create conditions table v1", NewAddTableMigration(conditionV1))

}

func addActionMigrations(mg *Migrator) {

	acionV1 := Table{
		Name: "actions",
		Columns: []*Column{
			{Name: "id", Type: DB_BigInt, IsPrimaryKey: true, IsAutoIncrement: true},
			{Name: "alarm_id", Type: DB_BigInt, Nullable: false},
			{Name: "throttle_period", Type: DB_NVarchar, Length: 100, Nullable: false},
			{Name: "notification_type", Type: DB_NVarchar, Length: 50, Nullable: false},
			{Name: "message", Type: DB_NVarchar, Length:  999, Nullable: false},
		},
		Indices: []*Index{
			{Cols: []string{"alarm_id"}, Type: IndexType},
		},
	}
	mg.AddMigration("create actions table v1", NewAddTableMigration(acionV1))

}

