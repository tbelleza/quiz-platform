# QuizLab — Test-Taking Platform with Analytics

QuizLab is a lightweight, full-stack web application that allows students to take quizzes and tests, submit their answers, and view their results. The platform also includes an analytics dashboard for visualizing student performance and identifying trends in assessment data.

The project demonstrates full-stack application development, cloud deployment, database management, and data analytics using AWS, Go, PostgreSQL, TypeScript, Vue.js, D3.js, and R.

## Features

- **Quiz Library:** Browse and select available quizzes and tests.
- **Test-Taking Interface:** Answer multiple-choice questions through a clean, responsive interface.
- **Automatic Grading:** Submit answers and receive immediate scores.
- **Results Tracking:** Store student quiz attempts and scores in a PostgreSQL database.
- **Analytics Dashboard:** Visualize average quiz scores and student performance using D3.js.
- **Statistical Analysis:** Analyze assessment data using R and data science libraries.

## Tech Stack

| Technology | Purpose |
|---|---|
| **Vue.js** | Frontend framework for building the user interface |
| **TypeScript** | Type-safe frontend development |
| **PrimeVue** | UI component library |
| **D3.js** | Interactive data visualizations |
| **Go (Golang)** | Backend REST API and automatic grading |
| **PostgreSQL** | Relational database for quizzes, questions, and results |
| **R** | Statistical analysis and reporting |
| **dplyr** | Data manipulation and aggregation |
| **AWS EC2** | Cloud hosting for backend services |
| **Docker** | Containerization of backend services and database |
| **GitHub Pages** | Hosting and deployment of the frontend |

## Architecture

The application follows a simple client-server architecture.

```text
             GitHub Pages
                  |
       Vue.js + TypeScript
         PrimeVue + D3.js
                  |
              REST API
                  |
             AWS EC2
                  |
             Go Backend
                  |
             PostgreSQL
                  |
             R Analytics
```

The Vue frontend communicates with a Go REST API hosted on AWS. The backend retrieves quiz questions, grades submissions, and stores results in PostgreSQL.

D3.js visualizes performance metrics in the frontend, while R provides additional statistical analysis and reporting.

## Project Structure

```text
quiz-platform/
├── frontend/          # Vue.js, TypeScript, PrimeVue, D3.js
├── backend/           # Go REST API
├── database/          # PostgreSQL schema and seed data
├── analytics/         # R scripts for statistical analysis
├── docker-compose.yml # Container configuration
└── README.md
```

## Analytics

QuizLab includes two approaches to analyzing assessment performance.

**D3.js Dashboard**

Visualizes average quiz scores and allows users to compare performance across assessments.

**R Statistical Analysis**

Uses R and `dplyr` to calculate:

- Average quiz scores
- Median scores
- Student pass rates
- Total quiz attempts

Statistical summaries can be exported as CSV files for further analysis.

## Deployment

The application uses a split deployment architecture:

- **Frontend:** GitHub Pages
- **Backend:** AWS EC2
- **Database:** PostgreSQL running in Docker on AWS EC2
- **Analytics:** R scripts for statistical reporting

The frontend communicates with the AWS-hosted backend through an HTTPS REST API.

## Project Goals

This project was developed as a minimal demonstration of:

- Building a full-stack web application
- Designing RESTful APIs using Go
- Managing relational data with PostgreSQL
- Creating responsive interfaces with Vue.js and PrimeVue
- Visualizing application data with D3.js
- Performing statistical analysis using R
- Deploying applications using AWS and GitHub Pages

## Future Improvements

Potential enhancements include student authentication, instructor dashboards, quiz creation tools, timed assessments, detailed answer reviews, and more advanced student performance analytics.

## License

This project is intended for educational and demonstration purposes.
