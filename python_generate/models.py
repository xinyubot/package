import sys

Import = """
import (
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)
"""

TName = """
// TName returns the actual table name for the receiver t. 
func (t *ttttttttt) TName() string {
    return "" + func() string {panic("Table Name undefined"); return ""}()
}
"""

InsertOne = """
// InsertOne inserts one row into the corresponding table. 
func (t *ttttttttt) InsertOne(db *gorm.DB) (err error) {
	if t == nil {
		return errors.New("empty entry")
	}

    // insert
	if err = db.Table(t.TName()).Create(t).Error; err != nil {
        return err
    }

    return nil
}
"""

InsertMulti = """
// InsertMulti inserts multiple rows into the corresponding table. 
func (t *ttttttttt) InsertMulti(db *gorm.DB, rows *[]ttttttttt) (err error) {
	if rows == nil {
		return errors.New("empty entry")
	}
    if len(*rows) == 0 {
        return errors.New("empty rows")
    }

    // insert
    if err = db.Table((*rows)[0].TName()).Create(rows).Error; err != nil {
        return err
    }

	return nil
}
"""

UpdateColumns = """
// UpdateColumns updates, if the row exists, the selected columns with new value specified by fields, and inserts if the row does not exist. 
func (t *ttttttttt) UpdateColumns(db *gorm.DB, fields map[string]any) (err error) {
	if t == nil {
		return errors.New("empty entry")
	}

    // check https://gorm.io/docs/create.html#Upsert-x2F-On-Conflict
	if err = db.Table(t.TName()).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: ""}},
		DoUpdates: clause.Assignments(fields),
	}).Create(t).Error; err != nil {
        return err
    }

    return nil
}
"""

FindOne = """
// FindOne finds the first entry that matches the given filtering options. 
func (t *ttttttttt) FindOne(db *gorm.DB) (ret ttttttttt, err error) {
	if t == nil {
		return ret, errors.New("empty entry")
	}

	// query
	if err = db.Table(t.TName()).First(&ret).Error; err != nil {
        return ret, err
    }

	return ret, nil
}
"""

FindAll = """
// FindAll finds all entries that match the given filtering options. 
func (t *ttttttttt) FindAll(db *gorm.DB) (ret []ttttttttt, count int, err error) {
	if t == nil {
		return ret, 0, errors.New("empty entry")
	}

	// query
	err = db.Table(t.TName()).Find(&ret).Error
	count = len(ret)

	return ret, count, err
}
"""

def main():
    try: 
        file_name = sys.argv[1]
        package_name = sys.argv[2]
        struct_name = "".join([string[0].upper() + string[1:] for string in file_name[0:-3].split("_")])
        import_str = "import ("
        d = {
            import_str: [Import, False],
            "TName() string": [TName, False],
            "InsertOne": [InsertOne, False],
            "InsertMulti": [InsertMulti, False],
            "UpdateColumns": [UpdateColumns, False],
            "FindOne": [FindOne, False],
            "FindAll": [FindAll, False]
        }
        registry = []
        lines = []

        with open(file_name, "r") as instream:
            for line in instream:
                for key, val in d.items():
                    if key in line:
                        val[1] = True
                lines.append(line)
            for key, val in d.items():
                func_str, ok = val
                if not ok:
                    registry.append(func_str)
                
        output = ""
        # old code
        for line in lines:
            output += line
            # import
            if "package " + package_name in line:
                func_str, ok = d[import_str]
                if not ok:
                    output += func_str
        # generated code
        output += "".join([func_str.replace("ttttttttt", struct_name) for func_str in registry if func_str != Import])

        with open(file_name, "w") as outstream:
            outstream.write(output)
            if len(registry) > 0:
                print("package {0} file {1}: orm struct {2} processed. ".format(package_name, file_name, struct_name))

    except Exception as e:
        print("package {0} file {1}: error: {2}. ".format(package_name, file_name, e))

    return  

if __name__ == '__main__':
    main()
