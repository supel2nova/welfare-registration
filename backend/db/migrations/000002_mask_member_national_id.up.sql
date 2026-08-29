UPDATE applications
SET applicant_snapshot = jsonb_set(
  applicant_snapshot,
  '{family,members}',
  (
    SELECT COALESCE(jsonb_agg(
      (m - 'national_id') || jsonb_build_object(
        'national_id_mask',
        CASE
          WHEN m->>'national_id' IS NULL OR length(m->>'national_id') <> 13 THEN 'null'::jsonb
          ELSE to_jsonb(
            left(m->>'national_id', 1) || '-' ||
            substring(m->>'national_id' from 2 for 4) || '-xxxxx-xx-' ||
            right(m->>'national_id', 1))
        END
      )
      ORDER BY ord
    ), '[]'::jsonb)
    FROM jsonb_array_elements(applicant_snapshot->'family'->'members') WITH ORDINALITY AS t(m, ord)
  )
)
WHERE jsonb_typeof(applicant_snapshot->'family'->'members') = 'array'
  AND EXISTS (
    SELECT 1 FROM jsonb_array_elements(applicant_snapshot->'family'->'members') e
    WHERE e ? 'national_id'
  );
