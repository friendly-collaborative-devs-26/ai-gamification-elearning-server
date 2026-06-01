# 🗄️ Documentación de Base de Datos

Este documento describe la estructura de la base de datos, sus tablas, relaciones y decisiones de diseño.

---

## 📐 Diagrama de Relaciones (ERD simplificado)

```
users ──────────────────────────────────────────────────────────┐
  │                                                              │
  ├──< accounts          (proveedores OAuth)                     │
  │                                                              │
  ├──< user_permissions >── permissions                          │
  │                                                              │
  ├──< block             (creator_id)  ──< block (parent_id)    │
  │         │                                                    │
  │         └──< exercises (creator_id / block_id) ─────────────┤
  │                   │                                          │
  │                   └──< exercise_status (user_id) ───────────┤
  │                              │                               │
  │                              └──< exercise_executions        │
  │                                                              │
  ├──< user_scores >── block                                     │
  │                                                              │
  └──< buyed_store_items >── store_items                         │
                                                                 │
  ──────────────────────────────────────────────────────────────┘
```

---

## 🔢 Enumeraciones

### `block_type`

Define la jerarquía de bloques de contenido.

| Valor      | Descripción                          |
| ---------- | ------------------------------------ |
| `road_map` | Nivel superior: agrupa varios cursos |
| `course`   | Nivel intermedio: agrupa módulos     |
| `module`   | Nivel inferior: contiene ejercicios  |

> **Jerarquía:** `road_map` → `course` → `module` → `exercises`

### `auth_provider`

Métodos de autenticación soportados.

| Valor    | Descripción                     |
| -------- | ------------------------------- |
| `local`  | Registro con email y contraseña |
| `google` | OAuth con Google                |
| `github` | OAuth con GitHub                |

---

## 📋 Tablas

### `users` — Usuarios del sistema

Tabla central del sistema. Almacena tanto usuarios con contraseña local como usuarios autenticados via OAuth.

| Columna      | Tipo      | Restricciones    | Descripción                           |
| ------------ | --------- | ---------------- | ------------------------------------- |
| `id`         | integer   | PK               | Identificador único                   |
| `username`   | varchar   | UNIQUE, NOT NULL | Nombre de usuario                     |
| `email`      | varchar   | UNIQUE, NOT NULL | Correo electrónico                    |
| `password`   | varchar   | nullable         | Hash de contraseña. Nulo si usa OAuth |
| `image`      | varchar   | nullable         | URL del avatar del usuario            |
| `created_at` | timestamp | default: `now()` | Fecha de registro                     |
| `updated_at` | timestamp | default: `now()` | Última modificación                   |

---

### `accounts` — Cuentas OAuth vinculadas

Permite que un usuario vincule múltiples proveedores OAuth a su cuenta.

| Columna               | Tipo      | Restricciones    | Descripción                        |
| --------------------- | --------- | ---------------- | ---------------------------------- |
| `user_id`             | int       | FK → `users.id`  | Usuario al que pertenece la cuenta |
| `provider`            | string    |                  | Ej: `"google"`, `"github"`         |
| `provider_account_id` | string    |                  | ID del usuario en el proveedor     |
| `created_at`          | timestamp | default: `now()` |                                    |
| `updated_at`          | timestamp | default: `now()` |                                    |

**Índices:**

- `user_id`
- `(provider, provider_account_id)` — UNIQUE _(evita duplicados por proveedor)_

---

### `permissions` — Catálogo de permisos

Define los roles y permisos disponibles en el sistema.

| Columna       | Tipo      | Restricciones    | Descripción                |
| ------------- | --------- | ---------------- | -------------------------- |
| `id`          | integer   | PK               |                            |
| `name`        | varchar   |                  | Ej: `"editor"`, `"solver"` |
| `description` | varchar   |                  | Descripción del permiso    |
| `created_at`  | timestamp | default: `now()` |                            |

---

### `user_permissions` — Asignación de permisos

Tabla de unión (muchos a muchos) entre usuarios y permisos.

| Columna         | Tipo      | Restricciones         | Descripción |
| --------------- | --------- | --------------------- | ----------- |
| `id`            | integer   | PK                    |             |
| `user_id`       | integer   | FK → `users.id`       |             |
| `permission_id` | integer   | FK → `permissions.id` |             |
| `created_at`    | timestamp | default: `now()`      |             |

---

### `block` — Bloques de contenido

Estructura jerárquica auto-referenciada que organiza el contenido en `road_map` → `course` → `module`.

| Columna       | Tipo         | Restricciones         | Descripción                                |
| ------------- | ------------ | --------------------- | ------------------------------------------ |
| `id`          | integer      | PK                    |                                            |
| `name`        | varchar      |                       | Nombre del bloque                          |
| `description` | text         |                       | Descripción                                |
| `type`        | `block_type` |                       | Tipo: `road_map`, `course` o `module`      |
| `parent_id`   | integer      | FK → `block.id`, null | Referencia al bloque padre (auto-relación) |
| `creator_id`  | integer      | FK → `users.id`       | Usuario que creó el bloque                 |
| `created_at`  | timestamp    | default: `now()`      |                                            |
| `updated_at`  | timestamp    | default: `now()`      |                                            |

> **Nota de diseño:** `parent_id` es `null` para los `road_map` (nivel raíz). Esta auto-relación permite la jerarquía sin tablas separadas.

---

### `exercises` — Ejercicios

Unidad mínima de aprendizaje. Siempre pertenece a un bloque de tipo `module`.

| Columna         | Tipo      | Restricciones    | Descripción                                 |
| --------------- | --------- | ---------------- | ------------------------------------------- |
| `id`            | integer   | PK               |                                             |
| `block_id`      | integer   | FK → `block.id`  | Módulo al que pertenece                     |
| `creator_id`    | integer   | FK → `users.id`  | Usuario que creó el ejercicio               |
| `title`         | varchar   |                  | Título del ejercicio                        |
| `description`   | text      |                  | Enunciado del ejercicio                     |
| `exercise_url`  | text      |                  | URL del repositorio o recurso del ejercicio |
| `configuration` | json      |                  | Configuración personalizable del ejercicio  |
| `difficulty`    | enum      |                  | Nivel de dificultad                         |
| `score`         | integer   |                  | Puntos que otorga al completarse            |
| `created_at`    | timestamp | default: `now()` |                                             |
| `updated_at`    | timestamp | default: `now()` |                                             |

---

### `exercise_status` — Estado de un ejercicio por usuario

Registra el estado actual de un usuario en un ejercicio específico (un registro por usuario/ejercicio).

| Columna             | Tipo      | Restricciones       | Descripción                                |
| ------------------- | --------- | ------------------- | ------------------------------------------ |
| `id`                | integer   | PK                  |                                            |
| `user_id`           | integer   | FK → `users.id`     |                                            |
| `exercise_id`       | integer   | FK → `exercises.id` |                                            |
| `exercise_fork_url` | varchar   |                     | URL del fork del usuario para el ejercicio |
| `veredict`          | enum      |                     | `approved` o `rejected`                    |
| `updated_at`        | timestamp | default: `now()`    |                                            |

---

### `exercise_executions` — Historial de ejecuciones

Registra cada intento/ejecución de un ejercicio. Permite auditoría y seguimiento de progreso.

| Columna              | Tipo      | Restricciones             | Descripción                            |
| -------------------- | --------- | ------------------------- | -------------------------------------- |
| `id`                 | integer   | PK                        |                                        |
| `exercise_status_id` | integer   | FK → `exercise_status.id` | Estado al que pertenece esta ejecución |
| `execution_time`     | integer   |                           | Duración de la ejecución (ms o seg)    |
| `status`             | enum      |                           | `in_progress`, `completed`, `failed`   |
| `created_at`         | timestamp | default: `now()`          |                                        |

---

### `user_scores` — Puntajes por curso

Acumula el puntaje total de un usuario dentro de un curso específico.

| Columna       | Tipo      | Restricciones    | Descripción                           |
| ------------- | --------- | ---------------- | ------------------------------------- |
| `id`          | integer   | PK               |                                       |
| `user_id`     | integer   | FK → `users.id`  |                                       |
| `block_id`    | integer   | FK → `block.id`  | Apunta a un bloque de tipo `course`   |
| `total_score` | integer   | default: `0`     | Suma de puntos acumulados en el curso |
| `updated_at`  | timestamp | default: `now()` |                                       |

---

### `store_items` — Ítems de la tienda

Catálogo de artículos disponibles para compra con los puntos acumulados.

| Columna         | Tipo      | Restricciones    | Descripción                     |
| --------------- | --------- | ---------------- | ------------------------------- |
| `id`            | integer   | PK               |                                 |
| `name`          | varchar   |                  | Nombre del ítem                 |
| `description`   | text      |                  |                                 |
| `price`         | integer   |                  | Costo en puntos                 |
| `item_type`     | varchar   |                  | Categoría del ítem              |
| `configuration` | json      |                  | Configuración/metadata del ítem |
| `is_active`     | boolean   | default: `true`  | Si está disponible en la tienda |
| `created_at`    | timestamp | default: `now()` |                                 |

---

### `buyed_store_items` — Compras realizadas

Registro histórico de ítems comprados por cada usuario.

| Columna        | Tipo      | Restricciones         | Descripción     |
| -------------- | --------- | --------------------- | --------------- |
| `id`           | integer   | PK                    |                 |
| `user_id`      | integer   | FK → `users.id`       |                 |
| `item_id`      | integer   | FK → `store_items.id` |                 |
| `purchased_at` | timestamp | default: `now()`      | Fecha de compra |

---

## 🔗 Mapa de Relaciones (Foreign Keys)

| Tabla origen          | Columna              | →   | Tabla destino     | Columna | Tipo de relación |
| --------------------- | -------------------- | --- | ----------------- | ------- | ---------------- |
| `accounts`            | `user_id`            | →   | `users`           | `id`    | N:1              |
| `user_permissions`    | `user_id`            | →   | `users`           | `id`    | N:1              |
| `user_permissions`    | `permission_id`      | →   | `permissions`     | `id`    | N:1              |
| `block`               | `creator_id`         | →   | `users`           | `id`    | N:1              |
| `block`               | `parent_id`          | →   | `block`           | `id`    | N:1 (auto-ref.)  |
| `exercises`           | `block_id`           | →   | `block`           | `id`    | N:1              |
| `exercises`           | `creator_id`         | →   | `users`           | `id`    | N:1              |
| `exercise_status`     | `user_id`            | →   | `users`           | `id`    | N:1              |
| `exercise_status`     | `exercise_id`        | →   | `exercises`       | `id`    | N:1              |
| `exercise_executions` | `exercise_status_id` | →   | `exercise_status` | `id`    | N:1              |
| `user_scores`         | `user_id`            | →   | `users`           | `id`    | N:1              |
| `user_scores`         | `block_id`           | →   | `block`           | `id`    | N:1              |
| `buyed_store_items`   | `user_id`            | →   | `users`           | `id`    | N:1              |
| `buyed_store_items`   | `item_id`            | →   | `store_items`     | `id`    | N:1              |

---

## 💡 Decisiones de Diseño Relevantes

### Autenticación híbrida

El campo `password` en `users` es nullable para soportar usuarios que se registran exclusivamente via OAuth (Google/GitHub). La tabla `accounts` permite vincular múltiples proveedores a un mismo usuario.

### Jerarquía de contenido con auto-referencia

La tabla `block` usa `parent_id` para modelar la jerarquía `road_map → course → module` sin necesidad de tablas separadas. Esto simplifica la estructura pero requiere validar el `type` en la lógica de negocio.

### Separación entre estado e historial de ejercicios

- `exercise_status`: estado **actual** del usuario en un ejercicio (1 registro por par usuario/ejercicio).
- `exercise_executions`: historial de **cada intento** de ejecución (N registros por estado).

Esta separación facilita consultar el veredicto final sin recorrer todo el historial.

### Puntajes por curso

`user_scores` apunta a un `block` de tipo `course`, acumulando el puntaje total de todos los ejercicios completados dentro de ese curso. Permite rankings y gamificación a nivel de curso.
