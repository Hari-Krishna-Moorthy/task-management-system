# Task Management System

## Overview

The Task Management Systerm allows users to manage task. 

## 1. Technologies Used

- **Programming Language**: Go (Golang)
- **Web Framework**: Fiber
- **Authentication**: JWT (JSON Web Tokens)
- **Database**: MongoDB Atlas (Managed NoSQL database)

## 2. Models

```mermaid
erDiagram
    USER {
      string ID PK "Primary key (UUID)"
      string Username "Username of the user"
      string Password "Hashed password"
      date CreatedAt "Date of account creation"
    }

    TASK {
      string ID PK "Primary key (UUID)"
      string Title "Title of the task"
      string Description "Detailed description of the task"
      string Status "Status of the task (e.g., pending, completed)"
      date DueDate "Due date for the task"
      string UserID FK "Foreign key referencing USER ID"
    }

    USER ||--o{ TASK : "owns"
```
