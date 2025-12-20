# Database Setup Instructions - Tibia Web Server

## Overview
Esta carpeta contiene los scripts SQL necesarios para preparar tu base de datos existente para la Tibia Web Server.

**Tu base de datos actual:**
- ✅ Schema.sql (v1.0) - Tablas principales
- ✅ z-001-migrate-v01-to-v02.sql - Migración aplicada
- ✅ z-002-migrate-v02-to-v03.sql - Migración aplicada

## Lo Que Necesitas Agregar

### 🔴 NUEVA TABLA REQUERIDA
**news** - Sistema de noticias para la página web
```sql
CREATE TABLE IF NOT EXISTS news (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title VARCHAR(255) NOT NULL,
	content TEXT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### 🟡 ACTUALIZACIÓN DE TABLA EXISTENTE
**Guilds** - Agregar columna description (si no existe)
```sql
ALTER TABLE Guilds ADD COLUMN description TEXT DEFAULT '';
```

## Tablas Que YA Tienes

✅ **Accounts** - Email, Auth, PremiumEnd, etc.
✅ **Characters** - AccountID, Name, Level, Profession, etc.
✅ **CharacterRights** - Para permisos (gamemaster, etc.)
✅ **Guilds** - WorldID, GuildID, Name, LeaderID, Created
✅ **GuildMembers** - GuildID, CharacterID, Rank, Title
✅ **GuildInvites** - GuildID, CharacterID, RecruiterID
✅ **GuildRanks** - Para rangos de guild
✅ **Houses** - HouseID, Name, Rent, Town, GuildHouse, etc.
✅ **HouseOwners** - OwnerID, PaidUntil
✅ **HouseAuctions** - BidderID, BidAmount, FinishTime
✅ **HouseTransfers** - Para transferencias de casas
✅ **HouseAssignments** - Para historial de asignaciones

## Cómo Ejecutar

### Opción 1: Ejecutar el script completo (recomendado)
```bash
sqlite3 tibia.db < database/migrations.sql
```

### Opción 2: Ejecutar comandos individuales
```bash
# Crear tabla news
sqlite3 tibia.db < database/migrations.sql

# Si la tabla Guilds ya existe, agregar columna description
sqlite3 tibia.db "ALTER TABLE Guilds ADD COLUMN description TEXT DEFAULT '';"
```

### Opción 3: Usar SQLite Browser
1. Abre tu base de datos en SQLite Browser
2. Abre la pestaña "Execute SQL"
3. Copia y pega el contenido de `database/migrations.sql`
4. Ejecuta

## Verificación Post-Instalación

Verifica que se creó la tabla news:
```bash
sqlite3 tibia.db "SELECT * FROM sqlite_master WHERE type='table' AND name='news';"
```

Verifica que Guilds tiene la columna description:
```bash
sqlite3 tibia.db "PRAGMA table_info(Guilds);"
```

## Columnas de Cada Tabla

### news (NUEVA)
| Columna | Tipo | Notas |
|---------|------|-------|
| id | INTEGER | Primary Key, Auto-increment |
| title | VARCHAR(255) | Título de la noticia |
| content | TEXT | Contenido completo |
| created_at | TIMESTAMP | Auto fecha de creación |

### Guilds (ACTUALIZADA)
Se agregará:
| Columna | Tipo | Notas |
|---------|------|-------|
| description | TEXT | Descripción de la guild (nueva) |

## Notas Importantes

⚠️ **CREATE TABLE IF NOT EXISTS** - El script es seguro ejecutar varias veces, no creará duplicados

⚠️ **ALTER TABLE** - Si la columna description ya existe, ignorará el comando

⚠️ **Índices** - Se crean automáticamente para mejor rendimiento en búsquedas

⚠️ **Timestamps** - En la tabla news, se usan TIMESTAMP con DEFAULT CURRENT_TIMESTAMP

## Soporte

Si encuentras errores:

1. **"Error: near "CREATE TABLE": syntax error"** 
   → Verifica que el archivo migrations.sql se descargó correctamente

2. **"Error: table "news" already exists"** 
   → La tabla ya existe, es seguro ignorar este error

3. **"Error: duplicate column name"** 
   → La columna description ya existe en Guilds, es seguro ignorar

4. **Problemas de permisos**
   → Asegúrate de tener acceso de lectura/escritura a tibia.db
