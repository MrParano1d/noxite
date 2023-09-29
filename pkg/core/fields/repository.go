package fields

import (
	"fmt"
)

type Repository struct {
	URL       RequiredString
	Type      *RequiredString
	Directory *RequiredString
}

// builder

type repositoryBuilder struct {
	url       string
	typ       string
	directory string
}

func RepositoryBuilder() *repositoryBuilder {
	return &repositoryBuilder{}
}

func (b *repositoryBuilder) URL(url string) *repositoryBuilder {
	b.url = url
	return b
}

func (b *repositoryBuilder) Type(typ string) *repositoryBuilder {
	b.typ = typ
	return b
}

func (b *repositoryBuilder) Directory(directory string) *repositoryBuilder {
	b.directory = directory
	return b
}

func (b *repositoryBuilder) Build() (Repository, error) {
	url, err := RequiredStringFromString(b.url)
	if err != nil {
		return Repository{}, err
	}

	repo := Repository{
		URL: url,
	}

	if b.typ != "" {
	typ, err := RequiredStringFromString(b.typ)
	if err != nil {
		return Repository{}, err
	}
	repo.Type = &typ
	}

	if b.directory != "" {
	directory, err := RequiredStringFromString(b.directory)
	if err != nil {
		return Repository{}, err
	}
	repo.Directory = &directory
	}

	return repo, nil
}

// errors

type RepositoryValidationError struct {
	Repository string
	Err        error
}

func (e *RepositoryValidationError) Error() string {
	return fmt.Sprintf("invalid repository: %s: %s", e.Repository, e.Err.Error())
}
