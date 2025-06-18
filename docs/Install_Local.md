# Local Installation Guide (Development/Testing)

This guide helps you run Fluently on your local machine for development or grading. No domain or SSL required. All services run on `localhost` using Docker Compose.

---

## Prerequisites
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [Git](https://git-scm.com/)

---

## 1. Clone the Repository
```bash
git clone https://github.com/IU-Capstone-Project-2025/Fluently.git
cd Fluently/backend
```

---

## 2. Create Environment File
Create a `.env` file in the `backend` directory:
```env
# JWT
JWT_SECRET=supersecretjwtkey
JWT_EXPIRATION=24h
REFRESH_EXPIRATION=720h

# App
APP_NAME=fluently
# Port for the backend API, recommended to use something other than 8080 to avoid conflicts
APP_PORT=8070 

# Database
DB_USER=postgres
DB_PASSWORD=securepassword
DB_HOST=postgres
DB_PORT=5432
DB_NAME=postgres

# Directus configuration
# Directus port recommended to use something other than 8080 to avoid conflicts
DIRECTUS_PORT=8055
DIRECTUS_ADMIN_EMAIL=admin@example.com
DIRECTUS_ADMIN_PASSWORD=admin
DIRECTUS_SECRET_KEY=supersecretkey

# Google OAuth
IOS_GOOGLE_CLIENT_ID=some-ios-client-id.apps.googleusercontent.com
ANDROID_GOOGLE_CLIENT_ID=some-android-client-id.apps.googleusercontent.com
WEB_GOOGLE_CLIENT_ID=some-web-client-id.apps.googleusercontent.com
```

---

## 3. Start All Services
From Fluently/backend:
```bash
docker compose up -d --build
```

---

## 4. Access Services
- **Backend API:** [http://localhost:8070/api/v1/](http://localhost:8080/api/v1/)
- **Swagger UI:** [http://localhost:8070/swagger/](http://localhost:8080/swagger/)
- **Directus:** [http://localhost:8055/admin](http://localhost:8055/)

---

## 5. Mobile Apps (Optional)
- **Android:** See [Android App README](android-app/README.md) for instructions.
- **iOS:** Open `ios-app/Fluently/Fluently.xcodeproj` in Xcode and run on a simulator. Set API base URL to `http://localhost:8080/`.

---

## 6. Stopping Services
From Fluently/backend:
```bash
docker compose down
```

---

## Notes
- No domain or SSL is required for local testing.
- For production setup, see [Full Installation Guide](Install_Full.md).
