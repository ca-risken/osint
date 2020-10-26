package osint

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Osint

// Validate ListOsintRequest
func (r *ListOsintRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
	)
}

// Validate GetOsintRequest
func (r *GetOsintRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.OsintId, validation.Required),
	)
}

// Validate PutOsintRequest
func (r *PutOsintRequest) Validate() error {
	if r.Osint == nil {
		return errors.New("Required Osint")
	}
	if err := validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required, validation.In(r.Osint.ProjectId)),
	); err != nil {
		return err
	}
	return r.Osint.Validate()
}

// Validate DeleteDataSourceRequest
func (r *DeleteOsintRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.OsintId, validation.Required),
	)
}

// OsintDataSource

// Validate ListOsintDataSourceRequest
func (r *ListOsintDataSourceRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
	)
}

// Validate GetOsintDataSourceRequest
func (r *GetOsintDataSourceRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.OsintDataSourceId, validation.Required),
	)
}

// Validate PutOsintDataSourceRequest
func (r *PutOsintDataSourceRequest) Validate() error {
	if r.OsintDataSource == nil {
		return errors.New("Required OsintDataSource")
	}
	if err := validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
	); err != nil {
		return err
	}
	return r.OsintDataSource.Validate()
}

// Validate DeleteDataSourceRequest
func (r *DeleteOsintDataSourceRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.OsintDataSourceId, validation.Required),
	)
}

// RelOsintDataSource

// Validate ListRelOsintDataSourceRequest
func (r *ListRelOsintDataSourceRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
	)
}

// Validate GetRelOsintDataSourceRequest
func (r *GetRelOsintDataSourceRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.RelOsintDataSourceId, validation.Required),
	)
}

// Validate PutRelOsintDataSourceRequest
func (r *PutRelOsintDataSourceRequest) Validate() error {
	if r.RelOsintDataSource == nil {
		return errors.New("Required RelOsintDataSource")
	}
	if err := validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required, validation.In(r.RelOsintDataSource.ProjectId), validation.Required),
	); err != nil {
		return err
	}
	return r.RelOsintDataSource.Validate()
}

// Validate DeleteResultRequest
func (r *DeleteRelOsintDataSourceRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.RelOsintDataSourceId, validation.Required),
	)
}

// Validate ListRelOsintDetectWordRequest
func (r *ListRelOsintDetectWordRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
	)
}

// Validate GetRelOsintDetectWordRequest
func (r *GetRelOsintDetectWordRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.RelOsintDetectWordId, validation.Required),
	)
}

// Validate PutRelOsintDetectWordRequest
func (r *PutRelOsintDetectWordRequest) Validate() error {
	if r.RelOsintDetectWord == nil {
		return errors.New("Required RelOsintDetectWord")
	}
	if err := validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required, validation.In(r.RelOsintDetectWord.ProjectId), validation.Required),
	); err != nil {
		return err
	}
	return r.RelOsintDetectWord.Validate()
}

// Validate DeleteResultRequest
func (r *DeleteRelOsintDetectWordRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.RelOsintDetectWordId, validation.Required),
	)
}

// Validate ListOsintDetectWordRequest
func (r *ListOsintDetectWordRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
	)
}

// Validate GetOsintDetectWordRequest
func (r *GetOsintDetectWordRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.OsintDetectWordId, validation.Required),
	)
}

// Validate PutOsintDetectWordRequest
func (r *PutOsintDetectWordRequest) Validate() error {
	if r.OsintDetectWord == nil {
		return errors.New("Required OsintDetectWord")
	}
	if err := validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required, validation.In(r.OsintDetectWord.ProjectId), validation.Required),
	); err != nil {
		return err
	}
	return r.OsintDetectWord.Validate()
}

// Validate DeleteResultRequest
func (r *DeleteOsintDetectWordRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.OsintDetectWordId, validation.Required),
	)
}

// Validate StartOsintRequest
func (r *StartOsintRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.RelOsintDataSourceId, validation.Required),
	)
}

/**
 * Entity
**/

// Validate Osint
func (d *OsintForUpsert) Validate() error {
	return validation.ValidateStruct(d,
		validation.Field(&d.ResourceType, validation.Required, validation.Length(0, 50)),
		validation.Field(&d.ResourceName, validation.Required, validation.Length(0, 200)),
		validation.Field(&d.ProjectId, validation.Required),
	)
}

// Validate OsintDataSourceForUpsert
func (d *OsintDataSourceForUpsert) Validate() error {
	return validation.ValidateStruct(d,
		validation.Field(&d.Name, validation.Required, validation.Length(0, 50)),
		validation.Field(&d.Description, validation.Required, validation.Length(0, 200)),
		validation.Field(&d.MaxScore, validation.Required, validation.Min(0.0), validation.Max(99999.0)),
	)
}

// Validate RelOsintDataSourceForUpsert
func (r *RelOsintDataSourceForUpsert) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.OsintDataSourceId, validation.Required),
		validation.Field(&r.OsintId, validation.Required),
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.ScanAt, validation.Min(0), validation.Max(253402268399)), //  1970-01-01T00:00:00 ~ 9999-12-31T23:59:59
		validation.Field(&r.StatusDetail, validation.Length(0, 255)),
	)
}

// Validate RelOsintDetectWordForUpsert
func (r *RelOsintDetectWordForUpsert) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.RelOsintDataSourceId, validation.Required),
		validation.Field(&r.OsintDetectWordId, validation.Required),
		validation.Field(&r.ProjectId, validation.Required),
	)
}

// Validate OsintDetectWordForUpsert
func (r *OsintDetectWordForUpsert) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ProjectId, validation.Required),
		validation.Field(&r.Word, validation.Required, validation.Length(0, 50)),
	)
}
