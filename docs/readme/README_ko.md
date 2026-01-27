# PieMDM - 엔터프라이즈 마스터 데이터 관리 시스템

[![CI](https://github.com/pieteams/piemdm/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/pieteams/piemdm/actions/workflows/ci.yml)

[English](../../README.md)
| [简体中文](README_zh-CN.md)
| [繁體中文](README_zh-TW.md)
| **한국어**
| [Русский](README_ru.md)
| [Tiếng Việt](README_vi.md)
| [日本語](README_ja.md)

PieMDM은 엔터프라이즈 데이터 거버넌스를 위해 설계된 강력하고 사용자 친화적인 오픈 소스 마스터 데이터 관리(MDM) 시스템입니다. Go 백엔드와 Vue.js 프론트엔드로 구축되어 포괄적인 데이터 관리, 거버넌스 및 통합 기능을 제공합니다.

**프로젝트 웹사이트**: https://pieteams.github.io/piemdm/

## 🚀 주요 기능

- 데이터 관리 및 통합
- 마스터 데이터 모델링
- 데이터 거버넌스
- 시스템 통합
- 액세스 제어
- 워크플로우 관리

## 📋 요구 사항

- Go 1.24.12 이상
- Node.js 20 이상
- MySQL 8.0 이상
- Redis 6 이상

## 🚀 빠른 시작

### 1. 저장소 클론

```bash
git clone https://github.com/pieteams/piemdm.git
cd piemdm
```

### 2. 백엔드 설정

```bash
cd backend
# 종속성 설치
go mod tidy

# 구성 파일 복사
cp config/local.yml.example config/local.yml
# config/local.yml을 편집하여 데이터베이스 및 기타 설정 구성

# 데이터베이스 마이그레이션 실행
go run cmd/migration/main.go

# 백엔드 서비스 시작
go run cmd/server/main.go
```

### 3. 프론트엔드 설정

```bash
cd frontend
# 종속성 설치
npm install
# 또는
pnpm install

# 개발 서버 시작
npm run dev
# 또는
pnpm dev
```

### 4. 애플리케이션 접속

- 프론트엔드: http://localhost:8081
- 백엔드 API: http://localhost:8787
- API 문서: http://localhost:8787/swagger/index.html

## 🐳 Docker 배포

### Docker Compose 사용

```bash
# 모든 서비스 시작
docker-compose -f deploy/docker-compose.yml up -d

# 로그 보기
docker-compose -f deploy/docker-compose.yml logs -f

# 서비스 중지
docker-compose -f deploy/docker-compose.yml down
```

### Docker 수동 빌드

```bash
# 백엔드 이미지 빌드
cd backend
docker build -t piemdm-api:latest -f scripts/build/Dockerfile .

# 프론트엔드 이미지 빌드
cd frontend
docker build -t piemdm-web:latest .
```

## 🔧 구성

### 환경 변수

루트 디렉토리에 `.env` 파일을 생성합니다.

```env
# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=piemdm

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=your_password

# JWT
JWT_SECRET=your_jwt_secret_key

# Application
APP_ENV=development
APP_PORT=8787
```

### 데이터베이스 설정

```sql
CREATE DATABASE piemdm CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

## 📖 API 문서

API 문서는 Swagger를 사용하여 자동으로 생성되며 다음 주소에서 액세스할 수 있습니다.

- 개발 환경: http://localhost:8787/swagger/index.html
- 운영 환경: https://your-domain.com/swagger/index.html

## 🧪 테스트

### 백엔드 테스트

```bash
cd backend
# 모든 테스트 실행
make test

# 커버리지와 함께 테스트 실행
go test -cover ./...

# Mock 생성
make mock
```

### 프론트엔드 테스트

```bash
cd frontend
# 단위 테스트 실행
pnpm test

# E2E 테스트 실행
npm run test:e2e
```

## 🚀 배포

### 프로덕션 빌드

```bash
# 백엔드 빌드
cd backend
make build

# 프론트엔드 빌드
cd frontend
pnpm build
```

### 환경별 구성

- 개발 환경: `config/local.yml`
- 운영 환경: `config/prod.yml`

## 🤝 기여하기

기여를 환영합니다! 자세한 내용은 [기여 가이드](CONTRIBUTING.md)를 참조하십시오.

### 개발 워크플로우

1. 저장소 포크 (Fork)
2. 기능 브랜치 생성
3. 변경 사항 커밋
4. 테스트 추가
5. Pull Request 제출

### 코드 표준

- Go 모범 사례 및 규칙 준수
- 프론트엔드 코드에 ESLint 및 Prettier 사용
- 포괄적인 테스트 작성
- 문서 업데이트

## 📄 라이선스

이 프로젝트는 MIT 라이선스에 따라 라이선스가 부여됩니다. 자세한 내용은 [LICENSE](LICENSE) 파일을 참조하십시오.

## 📞 지원

- 📧 이메일: [jasen215@gmail.com]
- 🐛 이슈: [GitHub Issues](https://github.com/pieteams/piemdm/issues)
- 💬 토론: [GitHub Discussions](https://github.com/pieteams/piemdm/discussions)

## 🙏 감사의 말

이 프로젝트를 가능하게 해준 모든 기여자 및 오픈 소스 커뮤니티에 감사드립니다.

---

**이 프로젝트가 도움이 되었다면 Star ⭐를 눌러주세요!**
