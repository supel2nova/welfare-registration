-- UUID คงที่ เพื่อให้เขียน test และ curl ซ้ำได้ (SPEC §3.3)

INSERT INTO organizations (id, code, name_th, org_type, province_code, district_code, scope_policy) VALUES
  ('11111111-1111-1111-1111-111111111111','KTB-0001','ธนาคารกรุงไทย สาขาสีลม','BANK_BRANCH',NULL,NULL,'OWN_UNIT'),
  ('22222222-2222-2222-2222-222222222222','BAAC-0007','ธ.ก.ส. สาขาแม่ริม','BANK_BRANCH',NULL,NULL,'OWN_UNIT'),
  ('33333333-3333-3333-3333-333333333333','GSB-0012','ธนาคารออมสิน สาขาเชียงใหม่','BANK_BRANCH',NULL,NULL,'OWN_UNIT'),
  ('44444444-4444-4444-4444-444444444444','DOPA-5001','ที่ว่าการอำเภอเมืองเชียงใหม่','DISTRICT_OFFICE','50','5001','DISTRICT'),
  ('55555555-5555-5555-5555-555555555555','CGD-50','สำนักงานคลังจังหวัดเชียงใหม่','TREASURY_PROVINCIAL','50',NULL,'PROVINCE'),
  ('66666666-6666-6666-6666-666666666666','PATTAYA-01','สำนักงานเมืองพัทยา','SPECIAL_AREA','20',NULL,'OWN_UNIT'),
  ('99999999-9999-9999-9999-999999999999','CENTRAL','กรมบัญชีกลาง','CENTRAL',NULL,NULL,'NATIONWIDE')
ON CONFLICT (code) DO NOTHING;

INSERT INTO users (id, username, role, organization_id) VALUES
  ('aaaaaaaa-0000-0000-0000-000000000001','somchai.ktb','REGISTRAR','11111111-1111-1111-1111-111111111111'),
  ('aaaaaaaa-0000-0000-0000-000000000002','somying.baac','REGISTRAR','22222222-2222-2222-2222-222222222222'),
  ('aaaaaaaa-0000-0000-0000-000000000003','admin','ADMIN','99999999-9999-9999-9999-999999999999')
ON CONFLICT (username) DO NOTHING;

INSERT INTO api_clients (id, client_id, name_th, organization_id) VALUES
  ('bbbbbbbb-0000-0000-0000-000000000001','ktb-core','ระบบสาขากรุงไทย','11111111-1111-1111-1111-111111111111')
ON CONFLICT (client_id) DO NOTHING;
