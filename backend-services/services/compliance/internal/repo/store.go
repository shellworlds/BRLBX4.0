package repo

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/shellworlds/BRLBX4.0/backend-services/pkg/pgxutil"
)

type ConsentRecord struct {
	ID            string    `json:"id"`
	Auth0Subject  string    `json:"auth0_subject"`
	PolicyVersion string    `json:"policy_version"`
	Accepted      bool      `json:"accepted"`
	RecordedAt    time.Time `json:"recorded_at"`
}

type AuditEntry struct {
	ID           int64           `json:"id"`
	ActorSubject *string         `json:"actor_subject,omitempty"`
	Action       string          `json:"action"`
	ResourceType *string         `json:"resource_type,omitempty"`
	ResourceID   *string         `json:"resource_id,omitempty"`
	Metadata     json.RawMessage `json:"metadata"`
	CreatedAt    time.Time       `json:"created_at"`
}

type ContactSubmission struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Company   *string   `json:"company,omitempty"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

type Store struct {
	DB pgxutil.Querier
}

func (s *Store) InsertConsent(ctx context.Context, subject, version string, accepted bool) error {
	const q = `INSERT INTO consent_records (auth0_subject, policy_version, accepted) VALUES ($1,$2,$3)`
	_, err := s.DB.Exec(ctx, q, subject, version, accepted)
	return err
}

func (s *Store) LatestConsent(ctx context.Context, subject string) (version string, accepted bool, ok bool, err error) {
	const q = `
SELECT policy_version, accepted FROM consent_records
WHERE auth0_subject=$1 ORDER BY recorded_at DESC LIMIT 1`
	err = s.DB.QueryRow(ctx, q, subject).Scan(&version, &accepted)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", false, false, nil
		}
		return "", false, false, err
	}
	return version, accepted, true, nil
}

func (s *Store) InsertAudit(ctx context.Context, actor, action, resType, resID string, meta json.RawMessage) error {
	if meta == nil || len(meta) == 0 {
		meta = []byte(`{}`)
	}
	const q = `INSERT INTO audit_log (actor_subject, action, resource_type, resource_id, metadata) VALUES ($1,$2,$3,$4,$5)`
	_, err := s.DB.Exec(ctx, q, nullIfEmpty(actor), action, nullIfEmptyPtr(resType), nullIfEmptyPtr(resID), meta)
	return err
}

func nullIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func nullIfEmptyPtr(s string) *string {
	return nullIfEmpty(s)
}

func (s *Store) ListAudit(ctx context.Context, limit int) ([]AuditEntry, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	const q = `
SELECT id, actor_subject, action, resource_type, resource_id, metadata, created_at
FROM audit_log ORDER BY id DESC LIMIT $1`
	rows, err := s.DB.Query(ctx, q, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []AuditEntry
	for rows.Next() {
		var e AuditEntry
		if err := rows.Scan(&e.ID, &e.ActorSubject, &e.Action, &e.ResourceType, &e.ResourceID, &e.Metadata, &e.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

func (s *Store) InsertContact(ctx context.Context, name, email, company, message string) (string, error) {
	const q = `
INSERT INTO contact_submissions (name, email, company, message)
VALUES ($1,$2,$3,$4) RETURNING id::text`
	var id string
	err := s.DB.QueryRow(ctx, q, name, email, nullStr(company), message).Scan(&id)
	return id, err
}

func nullStr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func (s *Store) RequestDeletion(ctx context.Context, subject string) error {
	const q = `INSERT INTO deletion_requests (auth0_subject) VALUES ($1)`
	_, err := s.DB.Exec(ctx, q, subject)
	return err
}
