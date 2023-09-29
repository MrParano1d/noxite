package fields

import "fmt"

type Directories struct {
	Bin *RequiredString
	Man *RequiredString
}

// converters

type directoriesBuilder struct {
	bin *string
	man *string
}

func DirectoriesBuilder() *directoriesBuilder {
	return &directoriesBuilder{}
}

func (d *directoriesBuilder) Bin(bin string) *directoriesBuilder {
	d.bin = &bin
	return d
}

func (d *directoriesBuilder) Man(man string) *directoriesBuilder {
	d.man = &man
	return d
}

func (d *directoriesBuilder) Build() (Directories, error) {

	directories := Directories{}

	if d.bin != nil {
		bin, err := RequiredStringFromString(*d.bin)
		if err != nil {
			return Directories{}, &InvalidBinError{Value: *d.bin, Err: err}
		}
		directories.Bin = &bin
	}

	if d.man != nil {
		man, err := RequiredStringFromString(*d.man)
		if err != nil {
			return Directories{}, &InvalidManError{Value: *d.man, Err: err}
		}
		directories.Man = &man
	}

	return directories, nil
}

// errors

type InvalidBinError struct {
	Value string
	Err   error
}

func (e *InvalidBinError) Error() string {
	return fmt.Sprintf("invalid bin: %s: %s", e.Value, e.Err.Error())
}

type InvalidManError struct {
	Value string
	Err   error
}

func (e *InvalidManError) Error() string {
	return fmt.Sprintf("invalid man: %s: %s", e.Value, e.Err.Error())
}
