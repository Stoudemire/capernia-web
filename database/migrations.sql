-- Tibia Web Server Database Migrations
-- Para base de datos con schema.sql + migraciones v0.2 y v0.3
-- Solo agrega las tablas y columnas faltantes para compatibilidad con web server

-- ============================================================================
-- NUEVA TABLA REQUERIDA: NEWS
-- ============================================================================
-- Sistema de noticias para la página web
CREATE TABLE IF NOT EXISTS news (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title VARCHAR(255) NOT NULL,
	content TEXT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_news_created ON news(created_at);


-- ============================================================================
-- ACTUALIZAR TABLA EXISTENTE: GUILDS
-- ============================================================================
-- Agregar columna de descripción si no existe
-- (Ejecutar manualmente si la tabla ya existe)
-- ALTER TABLE Guilds ADD COLUMN description TEXT DEFAULT '';


-- ============================================================================
-- TABLAS YA PRESENTES EN TU SCHEMA (No agregar):
-- ============================================================================
-- Las siguientes tablas ya existen en tu database y NO necesitan ser creadas:
--
-- ✓ Accounts         - Ya existe con PremiumEnd
-- ✓ Characters       - Ya existe con AccountID
-- ✓ CharacterRights  - Ya existe (CharacterID, Name)
-- ✓ Guilds           - Ya existe (necesita columna description - ver arriba)
-- ✓ GuildMembers     - Ya existe
-- ✓ GuildInvites     - Ya existe
-- ✓ GuildRanks       - Ya existe
-- ✓ Houses           - Ya existe con HouseID, Name, Rent, Town, etc.
-- ✓ HouseOwners      - Ya existe con OwnerID, PaidUntil
-- ✓ HouseAuctions    - Ya existe con BidderID, BidAmount, FinishTime
-- ✓ HouseTransfers   - Ya existe
-- ✓ HouseAssignments - Ya existe


-- ============================================================================
-- INSTRUCCIONES DE EJECUCIÓN
-- ============================================================================
-- 1. Para SQLite (desde terminal):
--    sqlite3 tu_tibia.db < database/migrations.sql
--
-- 2. Para agregar columna description a Guilds (si la tabla ya existe):
--    sqlite3 tu_tibia.db "ALTER TABLE Guilds ADD COLUMN description TEXT DEFAULT '';"
--
-- 3. Verificar que se creó la tabla news:
--    sqlite3 tu_tibia.db ".tables"
--    sqlite3 tu_tibia.db "SELECT * FROM news;"
