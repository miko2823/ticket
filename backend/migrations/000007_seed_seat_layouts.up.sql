-- Tokyo Dome: Sturdy Summer Fest 2026 (10 seats, 2 sections)
-- Row A (front, premium) arced near stage, Row B (back) arced behind
UPDATE events SET seat_layout = '{
  "canvas": {"width": 600, "height": 400},
  "stage": {"x": 175, "y": 20, "width": 250, "height": 40, "label": "STAGE"},
  "sections": [
    {"id": "A", "label": "A席 (Front)", "color": "#E65100"},
    {"id": "B", "label": "B席 (Back)",  "color": "#1565C0"}
  ],
  "seats": [
    {"label": "A-1", "section": "A", "x": 140, "y": 140, "r": 18},
    {"label": "A-2", "section": "A", "x": 210, "y": 120, "r": 18},
    {"label": "A-3", "section": "A", "x": 300, "y": 110, "r": 18},
    {"label": "A-4", "section": "A", "x": 390, "y": 120, "r": 18},
    {"label": "A-5", "section": "A", "x": 460, "y": 140, "r": 18},
    {"label": "B-1", "section": "B", "x": 120, "y": 250, "r": 18},
    {"label": "B-2", "section": "B", "x": 210, "y": 230, "r": 18},
    {"label": "B-3", "section": "B", "x": 300, "y": 220, "r": 18},
    {"label": "B-4", "section": "B", "x": 390, "y": 230, "r": 18},
    {"label": "B-5", "section": "B", "x": 480, "y": 250, "r": 18}
  ]
}' WHERE id = 'a1111111-1111-1111-1111-111111111111';

-- Blue Note Tokyo: Jazz Night in Shibuya (6 seats, 2 sections)
-- Intimate jazz club, S seats close to stage, A seats behind
UPDATE events SET seat_layout = '{
  "canvas": {"width": 400, "height": 350},
  "stage": {"x": 100, "y": 15, "width": 200, "height": 35, "label": "STAGE"},
  "sections": [
    {"id": "S", "label": "S席 (Premium)", "color": "#FFB300"},
    {"id": "A", "label": "A席 (Regular)", "color": "#43A047"}
  ],
  "seats": [
    {"label": "S-1", "section": "S", "x": 130, "y": 120, "r": 20},
    {"label": "S-2", "section": "S", "x": 200, "y": 110, "r": 20},
    {"label": "S-3", "section": "S", "x": 270, "y": 120, "r": 20},
    {"label": "A-1", "section": "A", "x": 120, "y": 220, "r": 20},
    {"label": "A-2", "section": "A", "x": 200, "y": 210, "r": 20},
    {"label": "A-3", "section": "A", "x": 280, "y": 220, "r": 20}
  ]
}' WHERE id = 'a2222222-2222-2222-2222-222222222222';

-- Suntory Hall: Classical Evening (8 seats, 2 sections)
-- S seats centered front row, A seats wider row behind
UPDATE events SET seat_layout = '{
  "canvas": {"width": 500, "height": 380},
  "stage": {"x": 125, "y": 15, "width": 250, "height": 35, "label": "STAGE"},
  "sections": [
    {"id": "S", "label": "S席 (Premium)", "color": "#FFB300"},
    {"id": "A", "label": "A席 (Regular)", "color": "#43A047"}
  ],
  "seats": [
    {"label": "S-1", "section": "S", "x": 155, "y": 130, "r": 18},
    {"label": "S-2", "section": "S", "x": 215, "y": 120, "r": 18},
    {"label": "S-3", "section": "S", "x": 285, "y": 120, "r": 18},
    {"label": "S-4", "section": "S", "x": 345, "y": 130, "r": 18},
    {"label": "A-1", "section": "A", "x": 110, "y": 240, "r": 18},
    {"label": "A-2", "section": "A", "x": 195, "y": 230, "r": 18},
    {"label": "A-3", "section": "A", "x": 305, "y": 230, "r": 18},
    {"label": "A-4", "section": "A", "x": 390, "y": 240, "r": 18}
  ]
}' WHERE id = 'a3333333-3333-3333-3333-333333333333';
