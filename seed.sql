PRAGMA foreign_keys = ON;

-- ============================
-- GENRES (5)
-- ============================
INSERT INTO Genre (name) VALUES
('Action'),
('Drama'),
('Comedy'),
('Sci-Fi'),
('Thriller');

-- ============================
-- MOVIES (20)
-- ============================
INSERT INTO Movie (title, releaseYear, duration) VALUES
('Edge of Tomorrow', 2014, 113),
('Interstellar', 2014, 169),
('The Dark Knight', 2008, 152),
('Inception', 2010, 148),
('The Martian', 2015, 144),
('The Hangover', 2009, 100),
('Superbad', 2007, 113),
('Arrival', 2016, 116),
('Joker', 2019, 122),
('Mad Max: Fury Road', 2015, 120),
('Shutter Island', 2010, 138),
('Parasite', 2019, 132),
('Guardians of the Galaxy', 2014, 121),
('Iron Man', 2008, 126),
('The Social Network', 2010, 120),
('Whiplash', 2014, 107),
('Her', 2013, 126),
('La La Land', 2016, 128),
('The Revenant', 2015, 156),
('Django Unchained', 2012, 165);

-- ============================
-- ACTORS (15)
-- ============================
INSERT INTO Actor (name, birthDate) VALUES
('Tom Cruise', '1962-07-03'),
('Matthew McConaughey', '1969-11-04'),
('Christian Bale', '1974-01-30'),
('Leonardo DiCaprio', '1974-11-11'),
('Matt Damon', '1970-10-08'),
('Bradley Cooper', '1975-01-05'),
('Jonah Hill', '1983-12-20'),
('Amy Adams', '1974-08-20'),
('Joaquin Phoenix', '1974-10-28'),
('Charlize Theron', '1975-08-07'),
('Mark Ruffalo', '1967-11-22'),
('Chris Pratt', '1979-06-21'),
('Robert Downey Jr.', '1965-04-04'),
('Jesse Eisenberg', '1983-10-05'),
('Miles Teller', '1987-02-20');

-- ============================
-- MOVIE → GENRE RELATIONSHIPS
-- ============================
-- Action
INSERT INTO Movie_Genres VALUES (1, 1);
INSERT INTO Movie_Genres VALUES (3, 1);
INSERT INTO Movie_Genres VALUES (10, 1);
INSERT INTO Movie_Genres VALUES (13, 1);
INSERT INTO Movie_Genres VALUES (14, 1);

-- Drama
INSERT INTO Movie_Genres VALUES (2, 2);
INSERT INTO Movie_Genres VALUES (9, 2);
INSERT INTO Movie_Genres VALUES (12, 2);
INSERT INTO Movie_Genres VALUES (15, 2);
INSERT INTO Movie_Genres VALUES (16, 2);
INSERT INTO Movie_Genres VALUES (19, 2);

-- Comedy
INSERT INTO Movie_Genres VALUES (6, 3);
INSERT INTO Movie_Genres VALUES (7, 3);
INSERT INTO Movie_Genres VALUES (18, 3);

-- Sci-Fi
INSERT INTO Movie_Genres VALUES (1, 4);
INSERT INTO Movie_Genres VALUES (2, 4);
INSERT INTO Movie_Genres VALUES (4, 4);
INSERT INTO Movie_Genres VALUES (5, 4);
INSERT INTO Movie_Genres VALUES (8, 4);
INSERT INTO Movie_Genres VALUES (17, 4);

-- Thriller
INSERT INTO Movie_Genres VALUES (11, 5);
INSERT INTO Movie_Genres VALUES (9, 5);
INSERT INTO Movie_Genres VALUES (19, 5);
INSERT INTO Movie_Genres VALUES (20, 5);

-- Multi‑genre examples
INSERT INTO Movie_Genres VALUES (3, 5);  -- Dark Knight = Action + Thriller
INSERT INTO Movie_Genres VALUES (10, 5); -- Fury Road = Action + Thriller
INSERT INTO Movie_Genres VALUES (12, 5); -- Parasite = Drama + Thriller

-- ============================
-- MOVIE → ACTOR RELATIONSHIPS
-- ============================

-- Edge of Tomorrow
INSERT INTO Movie_Actors VALUES (1, 1);
INSERT INTO Movie_Actors VALUES (1, 10);

-- Interstellar
INSERT INTO Movie_Actors VALUES (2, 2);
INSERT INTO Movie_Actors VALUES (2, 8);

-- The Dark Knight
INSERT INTO Movie_Actors VALUES (3, 3);
INSERT INTO Movie_Actors VALUES (3, 9);

-- Inception
INSERT INTO Movie_Actors VALUES (4, 4);
INSERT INTO Movie_Actors VALUES (4, 8);

-- The Martian
INSERT INTO Movie_Actors VALUES (5, 5);
INSERT INTO Movie_Actors VALUES (5, 8);

-- The Hangover
INSERT INTO Movie_Actors VALUES (6, 6);
INSERT INTO Movie_Actors VALUES (6, 7);

-- Superbad
INSERT INTO Movie_Actors VALUES (7, 7);
INSERT INTO Movie_Actors VALUES (7, 15);

-- Arrival
INSERT INTO Movie_Actors VALUES (8, 8);
INSERT INTO Movie_Actors VALUES (8, 11);

-- Joker
INSERT INTO Movie_Actors VALUES (9, 9);

-- Mad Max: Fury Road
INSERT INTO Movie_Actors VALUES (10, 10);

-- Shutter Island
INSERT INTO Movie_Actors VALUES (11, 4);
INSERT INTO Movie_Actors VALUES (11, 11);

-- Parasite
INSERT INTO Movie_Actors VALUES (12, 11);

-- Guardians of the Galaxy
INSERT INTO Movie_Actors VALUES (13, 12);

-- Iron Man
INSERT INTO Movie_Actors VALUES (14, 13);

-- The Social Network
INSERT INTO Movie_Actors VALUES (15, 14);

-- Whiplash
INSERT INTO Movie_Actors VALUES (16, 15);

-- Her
INSERT INTO Movie_Actors VALUES (17, 9);
INSERT INTO Movie_Actors VALUES (17, 8);

-- La La Land
INSERT INTO Movie_Actors VALUES (18, 8);
INSERT INTO Movie_Actors VALUES (18, 15);

-- The Revenant
INSERT INTO Movie_Actors VALUES (19, 4);

-- Django Unchained
INSERT INTO Movie_Actors VALUES (20, 4);
INSERT INTO Movie_Actors VALUES (20, 5);
