
CREATE TABLE accountdb.users (
  id VARCHAR(255) NOT NULL PRIMARY KEY COMMENT 'ユーザーID',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = 'ユーザー';

CREATE TABLE accountdb.user_profiles (
  user_id VARCHAR(255) NOT NULL PRIMARY KEY COMMENT 'ユーザーID',
  email VARCHAR (255) NOT NULL UNIQUE COMMENT 'Eメールアドレス',
  display_name VARCHAR(20) NOT NULL COMMENT 'ユーザー名',
  screen_img_url VARCHAR(255) COMMENT 'スクリーン画像',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  CONSTRAINT user_profiles_ibfk_1 FOREIGN KEY(user_id) REFERENCES users(id) ON UPDATE RESTRICT ON DELETE RESTRICT
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = 'ユーザー';

CREATE TABLE accountdb.email_auth (
  user_id VARCHAR(255) NOT NULL PRIMARY KEY COMMENT 'ユーザーID',
  email VARCHAR(255) NOT NULL UNIQUE COMMENT 'Eメールアドレス',
  password_digest VARCHAR(255) NOT NULL COMMENT 'パスワードハッシュ',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  CONSTRAINT email_auth_ibfk_1 FOREIGN KEY(user_id) REFERENCES users(id) ON UPDATE RESTRICT ON DELETE RESTRICT
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = 'Eメール認証';

CREATE TABLE accountdb.email_confirmations (
  email VARCHAR(255) NOT NULL PRIMARY KEY COMMENT 'Eメールアドレス',
  confirmation_code VARCHAR(6) NOT NULL UNIQUE COMMENT 'メールアドレス確認番号',
  confirmed_at DATETIME,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE confirmation_code_idx (email, confirmation_code)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = 'Eメール認証登録';

CREATE TABLE accountdb.google_auth (
  user_id VARCHAR(255) NOT NULL PRIMARY KEY COMMENT 'ユーザーID',
  google_id VARCHAR(255) NOT NULL UNIQUE COMMENT 'Google Account ID',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  CONSTRAINT google_auth_ibfk_1 FOREIGN KEY(user_id) REFERENCES users(id) ON UPDATE RESTRICT ON DELETE RESTRICT
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = 'Google認証';

-- The schema must be the same as `src/internal/pubsubevent/pubsubevent.go`
CREATE TABLE accountdb.pubsub_events (
  id VARCHAR(255) NOT NULL PRIMARY KEY COMMENT 'イベントID',
  deduplicate_key VARCHAR(255) COMMENT '重複排除キー',
  topic VARCHAR(255) NOT NULL COMMENT 'トピック名',
  data TEXT NOT NULL COMMENT 'データ',
  is_published BOOLEAN NOT NULL DEFAULT FALSE COMMENT 'Publish済みか否か',
  published_at DATETIME COMMENT 'Publishした日時',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  INDEX is_published_idx (is_published)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = 'PubSubイベント';
