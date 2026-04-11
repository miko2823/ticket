ALTER TABLE events ADD COLUMN ticketing_starts_at TIMESTAMPTZ;
ALTER TABLE events ADD COLUMN ticketing_ends_at TIMESTAMPTZ;

-- Update seed data: ticketing opens 2 weeks before the event, closes 1 day before
UPDATE events SET
  ticketing_starts_at = starts_at - INTERVAL '2 weeks',
  ticketing_ends_at = starts_at - INTERVAL '1 day';

ALTER TABLE events ALTER COLUMN ticketing_starts_at SET NOT NULL;
ALTER TABLE events ALTER COLUMN ticketing_ends_at SET NOT NULL;
