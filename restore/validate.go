package restore

import (
	"fmt"
	"strings"

	"github.com/greenplum-db/gp-common-go-libs/dbconn"
	"github.com/greenplum-db/gp-common-go-libs/gplog"
	"github.com/greenplum-db/gpbackup/utils"
	"github.com/pkg/errors"
)

/*
 * This file contains functions related to validating user input.
 */

func validateFilterListsInBackupSet() {
	ValidateFilterSchemasInBackupSet(includeSchemas)
	ValidateFilterTablesInBackupSet(includeTables)
}

func ValidateFilterSchemasInBackupSet(schemaList utils.ArrayFlags) {
	schemaMap := make(map[string]bool, len(schemaList))
	for _, schema := range schemaList {
		schemaMap[schema] = true
	}
	if len(schemaList) > 0 {
		if !backupConfig.DataOnly {
			for _, entry := range globalTOC.PredataEntries {
				if _, ok := schemaMap[entry.Schema]; ok {
					delete(schemaMap, entry.Schema)
				}
				if len(schemaMap) == 0 {
					return
				}
			}
		} else {
			for _, entry := range globalTOC.DataEntries {
				if _, ok := schemaMap[entry.Schema]; ok {
					delete(schemaMap, entry.Schema)
				}
				if len(schemaMap) == 0 {
					return
				}
			}
		}
	} else {
		return
	}

	keys := make([]string, len(schemaMap))
	i := 0
	for k := range schemaMap {
		keys[i] = k
		i++
	}
	gplog.Fatal(errors.Errorf("Could not find the following schema(s) in the backup set: %s", strings.Join(keys, ", ")), "")
}

func ValidateFilterTablesInRestoreDatabase(connection *dbconn.DBConn, tableList utils.ArrayFlags) {
	if len(tableList) > 0 {
		utils.ValidateFQNs(tableList)
		quotedTablesStr := utils.SliceToQuotedString(tableList)
		query := fmt.Sprintf(`
SELECT
	n.nspname || '.' || c.relname AS string
FROM pg_namespace n
JOIN pg_class c ON n.oid = c.relnamespace
WHERE quote_ident(n.nspname) || '.' || quote_ident(c.relname) IN (%s)`, quotedTablesStr)
		resultTables := dbconn.MustSelectStringSlice(connection, query)
		if len(resultTables) > 0 {
			gplog.Fatal(nil, "Table %s already exists", resultTables[0])
		}
	}
}

func ValidateFilterTablesInBackupSet(tableList utils.ArrayFlags) {
	tableMap := make(map[string]bool, len(tableList))
	for _, table := range tableList {
		tableMap[table] = true
	}
	if len(tableList) > 0 {
		if !backupConfig.DataOnly {
			for _, entry := range globalTOC.PredataEntries {
				if entry.ObjectType != "TABLE" {
					continue
				}
				fqn := utils.MakeFQN(entry.Schema, entry.Name)
				if _, ok := tableMap[fqn]; ok {
					delete(tableMap, fqn)
				}
				if len(tableMap) == 0 {
					return
				}
			}
		} else {
			for _, entry := range globalTOC.DataEntries {
				fqn := utils.MakeFQN(entry.Schema, entry.Name)
				if _, ok := tableMap[fqn]; ok {
					delete(tableMap, fqn)
				}
				if len(tableMap) == 0 {
					return
				}
			}
		}
	} else {
		return
	}

	keys := make([]string, len(tableMap))
	i := 0
	for k := range tableMap {
		keys[i] = k
		i++
	}
	gplog.Fatal(errors.Errorf("Could not find the following table(s) in the backup set: %s", strings.Join(keys, ", ")), "")
}

func ValidateBackupFlagCombinations() {
	if backupConfig.SingleDataFile {
		if *numJobs != 1 {
			gplog.Fatal(errors.Errorf("Cannot use jobs flag when restoring backups with a single data file per segment."), "")
		}
	}
	if backupConfig.IncludeTableFiltered {
		if *createDB {
			gplog.Fatal(errors.Errorf("CREATE DATABASE statements are not included in table-filtered backups. Create the database manually before restoring."), "")
		} else if *restoreGlobals {
			gplog.Fatal(errors.Errorf("Global metadata is not backed up in table-filtered backup."), "")
		}
	}
}

func ValidateFlagCombinations() {
	utils.CheckMandatoryFlags("timestamp")
	utils.CheckExclusiveFlags("debug", "quiet", "verbose")
	utils.CheckExclusiveFlags("include-schema", "include-table", "include-table-file")
	utils.CheckExclusiveFlags("exclude-schema", "include-schema")
	utils.CheckExclusiveFlags("exclude-schema", "exclude-table", "include-table", "exclude-table-file", "include-table-file")
	utils.CheckExclusiveFlags("exclude-table", "exclude-table-file", "leaf-partition-data")
}
