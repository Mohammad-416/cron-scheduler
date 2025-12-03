cron-scheduler

cron-scheduler is a lightweight, interval-based job scheduler written in Go.
It works similarly to cron, but with one key difference: it remembers timer progress across restarts.

If a job was halfway through its interval when you shut the scheduler down, it resumes from the remaining time when you start it again.
If the downtime is longer than the interval, the job runs immediately and then falls back into its original rhythm.

This makes the scheduler predictable and reliable even when the process stops unexpectedly.

Why this exists

Traditional cron only runs jobs at fixed times.
Most interval schedulers ignore downtime and restart all timers from zero.

cron-scheduler fills the gap:

It remembers when each job last ran

It calculates how much of the interval was completed

It resumes the leftover delay on restart

If the interval is fully missed, it runs once immediately

After that, it returns to the original interval pattern

This behavior is ideal for periodic tasks that shouldn't drift over time.

Features

Simple YAML configuration for defining jobs

Each job runs in its own goroutine

Clean concurrency model with channels and context

State persistence through state.json

Resume timers from where they were paused

Automatic catch-up run when overdue

Graceful shutdown via OS signals

How it works

For each job, the scheduler tracks:

the last time it ran

how long the interval is

how much time passed before shutdown

On startup, it compares the last run time with the current time:

If elapsed < interval
→ resume from the remaining time

If elapsed ≥ interval
→ run once immediately
→ next run happens after interval - (elapsed % interval)

This keeps the scheduler aligned with the original timing.

Configuration

Create a config.yaml file:

jobs:
  - name: "say_hello"
    every: "10s"
    command: "echo 'Hello from cron-scheduler'"

  - name: "show_time"
    every: "15s"
    command: "date"


Each job contains:

name: unique identifier

every: interval duration (Go style, e.g., 5s, 10m, 1h)

command: shell command to execute

Installation

Build the scheduler:

go build -o cron-scheduler ./cmd


Run it:

./cron-scheduler -config=config.yaml -state=state.json


The first run will create state.json automatically.

Stopping & resuming

Stop the scheduler with:

CTRL + C


On shutdown, it saves:

last run timestamps

run counters

When you start it again, all jobs pick up exactly where they were.
