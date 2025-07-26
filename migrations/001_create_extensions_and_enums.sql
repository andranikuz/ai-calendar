-- Migration 001: Create extensions and enum types
-- This creates the foundation for the database

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create enum types
CREATE TYPE goal_category AS ENUM ('health', 'career', 'education', 'personal', 'financial', 'relationship');
CREATE TYPE goal_status AS ENUM ('draft', 'active', 'paused', 'completed', 'cancelled');
CREATE TYPE task_status AS ENUM ('pending', 'in_progress', 'completed', 'cancelled');
CREATE TYPE event_status AS ENUM ('tentative', 'confirmed', 'cancelled');
CREATE TYPE priority_level AS ENUM ('low', 'medium', 'high', 'critical');