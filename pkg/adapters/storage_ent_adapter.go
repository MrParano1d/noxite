package adapters

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/mrparano1d/noxite/ent"
	"github.com/mrparano1d/noxite/ent/repopackage"
	"github.com/mrparano1d/noxite/ent/version"
	"github.com/mrparano1d/noxite/pkg/core/entities"
	"github.com/mrparano1d/noxite/pkg/core/fields"
	"github.com/mrparano1d/noxite/pkg/core/ports"
)

type StorageEntAdapter struct {
	entClient *ent.Client
}

var _ ports.StoragePort = (*StorageEntAdapter)(nil)

func NewStorageEntAdapter(entClient *ent.Client) *StorageEntAdapter {
	return &StorageEntAdapter{
		entClient: entClient,
	}
}

func (s *StorageEntAdapter) createPackage(ctx context.Context, creatorID fields.EntityID, manifest *entities.PackageVersion) (*ent.RepoPackage, error) {
	return s.entClient.RepoPackage.Create().SetName(manifest.Name.String()).SetCreatorID(creatorID.Int()).Save(ctx)
}

func (s *StorageEntAdapter) reactivatePackage(ctx context.Context, creatorID fields.EntityID, pkg *ent.RepoPackage) error {
	return s.entClient.RepoPackage.Update().SetCreatorID(creatorID.Int()).SetNillableDeletedAt(nil).Where(repopackage.IDEQ(pkg.ID)).Exec(ctx)
}

func isVersionNewer(latest, newVersion string) (bool, error) {
	latestParts := strings.Split(latest, ".")
	newParts := strings.Split(newVersion, ".")

	isNewer := true
	for i := 0; i < len(latestParts); i++ {
		latestPart, err := strconv.Atoi(latestParts[i])
		if err != nil {
			return false, err
		}

		newPart, err := strconv.Atoi(newParts[i])
		if err != nil {
			return false, err
		}

		if newPart < latestPart {
			isNewer = false
			break
		}
	}

	return isNewer, nil
}

func (s *StorageEntAdapter) isVersionNewer(ctx context.Context, pkg *ent.RepoPackage, newVersion fields.RequiredString) (bool, error) {
	latestVersions, err := s.entClient.Version.Query().Where(version.PackageIDEQ(pkg.ID)).All(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return true, nil
		}
		return false, err
	}

	if len(latestVersions) == 0 {
		return true, nil
	}

	latest := latestVersions[0].Version

	for i := 1; i < len(latestVersions); i++ {
		v := latestVersions[i]
		if ok, err := isVersionNewer(latest, v.Version); err != nil {
			return false, err
		} else if ok {
			latest = v.Version
		}
	}

	return isVersionNewer(latest, newVersion.String())
}

func (s *StorageEntAdapter) createVersion(ctx context.Context, publisherID fields.EntityID, pkg *ent.RepoPackage, manifest *entities.PackageVersion) (*ent.Version, error) {
	query := s.entClient.Version.Create().
		SetVersion(manifest.Version.String()).
		SetNillableDescription(manifest.Description).
		SetKeywords(fields.StringsFromRequiredStrings(manifest.Keywords)).
		SetNillableLicense(manifest.License).
		SetContributors(manifest.Contributors).
		SetFunding(manifest.Funding).
		SetFiles(fields.StringsFromRequiredStrings(manifest.Files))

	if manifest.Homepage != nil {
		query = query.SetHomepage(manifest.Homepage.String())
	}

	if manifest.Bugs != nil {
		query = query.SetBugs(manifest.Bugs)
	}

	if manifest.Author != nil {
		query = query.SetAuthor(manifest.Author)
	}

	if manifest.Main != nil {
		query = query.SetMain(manifest.Main.String())
	}

	if manifest.Browser != nil {
		query = query.SetBrowser(manifest.Browser.String())
	}
	return query.SetBin(manifest.Bin).
		SetMan(manifest.Man).
		SetDirectories(manifest.Directories).
		SetRepository(manifest.Repository).
		SetScripts(manifest.Scripts).
		SetConfig(manifest.Config).
		SetDependencies(manifest.Dependencies).
		SetDevDependencies(manifest.DevDependencies).
		SetPeerDependencies(manifest.PeerDependencies).
		SetPeerDependenciesMeta(manifest.PeerDependenciesMeta).
		SetBundledDependencies(fields.StringsFromRequiredStrings(manifest.BundledDependencies)).
		SetOptionalDependencies(manifest.OptionalDependencies).
		SetOverrides(manifest.Overrides).
		SetEngines(manifest.Engines).
		SetOs(fields.StringsFromRequiredStrings(manifest.OS)).
		SetCPU(fields.StringsFromRequiredStrings(manifest.CPU)).
		SetNillablePrivate(manifest.Private).
		SetPublishConfig(manifest.PublishConfig).
		SetWorkspaces(fields.StringsFromRequiredStrings(manifest.Workspaces)).
		SetPackageID(pkg.ID).
		SetVersion(manifest.Version.String()).
		SetContentType(manifest.ContentType.String()).
		SetIntegrity(manifest.Integrity.String()).
		SetShasum(manifest.SHASUM.String()).
		SetLength(manifest.Length).
		SetData(manifest.Data.String()).
		SetPublisherID(publisherID.Int()).
		Save(ctx)
}

func (s *StorageEntAdapter) PublishPackage(ctx context.Context, creatorID fields.EntityID, manifest *entities.PackageVersion) error {
	var pkg *ent.RepoPackage
	var err error

	// check if package already exists

	pkg, err = s.entClient.RepoPackage.Query().Where(repopackage.NameEQ(manifest.Name.String())).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			pkg, err = s.createPackage(ctx, creatorID, manifest)
			if err != nil {
				return &ports.StorageAdapterPublishPackageError{
					Err: fmt.Errorf("failed to create package: %w", err),
				}
			}
		} else {
			return &ports.StorageAdapterPublishPackageError{
				Err: fmt.Errorf("failed to query package: %w", err),
			}
		}
	}

	if pkg.DeletedAt != nil {
		if err := s.reactivatePackage(ctx, creatorID, pkg); err != nil {
			return &ports.StorageAdapterPublishPackageError{
				Err: fmt.Errorf("failed to reactivate package: %w", err),
			}
		}
	}

	// if it does, check if the version is newer

	newer, err := s.isVersionNewer(ctx, pkg, manifest.Version)
	if err != nil {
		return &ports.StorageAdapterPublishPackageError{
			Err: fmt.Errorf("failed to check if version is newer: %w", err),
		}
	}

	if !newer {
		return &ports.StorageAdapterPublishPackageError{
			Err: fmt.Errorf("package version is not newer"),
		}
	}

	// if it is, create a new version

	if _, err := s.createVersion(ctx, creatorID, pkg, manifest); err != nil {
		return &ports.StorageAdapterPublishPackageError{
			Err: fmt.Errorf("failed to create version: %w", err),
		}
	}

	return nil
}

func (s *StorageEntAdapter) GetPackage(ctx context.Context, name fields.PackageName, rev fields.RequiredString) (*entities.PackageVersion, error) {

	pkg, err := s.entClient.RepoPackage.Query().WithVersions(func(vq *ent.VersionQuery) {
		vq.Order(ent.Desc(version.FieldVersion))
	}).Where(repopackage.NameEQ(name.String()), repopackage.DeletedAtIsNil()).First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, &ports.StorageAdapterPackageNotFoundError{
				Name:    name,
				Version: rev,
			}
		}
		return nil, &ports.StorageAdapterGetPackageError{
			Name:    name,
			Version: rev,
			Err:     fmt.Errorf("failed to query package: %w", err),
		}
	}

	if rev.String() == "latest" {
		return packageVersionFromEntVersion(name.String(), pkg.Edges.Versions[0])
	} else {
		for _, v := range pkg.Edges.Versions {
			if v.Version == rev.String() {
				return packageVersionFromEntVersion(name.String(), v)
			}
		}
	}

	return nil, &ports.StorageAdapterPackageNotFoundError{
		Name:    name,
		Version: rev,
	}

}
