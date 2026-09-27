import { get, post, pick, think, safeJson, asArray, students, admins } from "./lib.js";

// Student session: log in once (handled in lib), browse the dashboard, drill into
// a course -> lesson -> topics/assignments/quizzes, and occasionally write
// (mark progress, submit an assignment). Mirrors what a real student does.
export function student() {
  const user = students[__VU % students.length];
  if (!user) return;

  // Landing: my courses + the things a student glances at on login.
  const coursesRes = get(user, "/dashboard/student", "GET /dashboard/student");
  get(user, "/announcement/s", "GET /announcement/s");
  get(user, "/timetable/s/timetable", "GET /timetable/s/timetable");
  get(user, "/attendance/s/me/summary", "GET /attendance/s/me/summary");
  think(1, 4);

  const courses = asArray(safeJson(coursesRes));
  if (courses.length) {
    const course = pick(courses);
    const lessonsRes = get(user, `/dashboard/student/${course.id}/lessons`, "GET /dashboard/student/{courseId}/lessons");
    think(1, 3);

    const lessons = asArray(safeJson(lessonsRes));
    if (lessons.length) {
      const lesson = pick(lessons);
      get(user, `/topic/s/lesson/${lesson.id}`, "GET /topic/s/lesson/{lessonId}");
      const assignmentsRes = get(user, `/assignment/s/lessons/${lesson.id}/assignments`, "GET /assignment/s/lessons/{lessonId}/assignments");
      get(user, `/quiz/s/lesson/${lesson.id}`, "GET /quiz/s/lesson/{lessonId}");
      think(2, 5);

      // ~30% of sessions mark a lesson complete.
      if (Math.random() < 0.3) {
        post(user, `/dashboard/student/${lesson.id}/progress`, { completed: true }, "POST /dashboard/student/{lessonId}/progress");
      }
      // ~20% submit an assignment (excludes the Lambda code-runner on purpose).
      const assignments = asArray(safeJson(assignmentsRes));
      if (assignments.length && Math.random() < 0.2) {
        const asg = pick(assignments);
        post(user, `/assignment/s/assignments/${asg.id}/submit`, { code: "print('load test submission')" }, "POST /assignment/s/assignments/{assignmentId}/submit");
      }
    }
  }

  // Wrap up: check grades/submissions.
  get(user, "/assignment/s/submissions/summary", "GET /assignment/s/submissions/summary");
  think(1, 3);
}

// Admin session: dashboards and management lists. Reads only by design — admin
// writes create real rows (students, courses) that need unique data and would
// pollute the DB, so we don't fuzz them under load here.
export function admin() {
  const user = admins[__VU % admins.length];
  if (!user) return;

  get(user, "/dashboard", "GET /dashboard");
  think(1, 2);
  get(user, "/dashboard/admin?page=1&pageSize=20", "GET /dashboard/admin");
  get(user, "/course", "GET /course");
  const batchesRes = get(user, "/batch", "GET /batch");
  get(user, "/program/a", "GET /program/a");
  get(user, "/announcement/a", "GET /announcement/a");
  think(1, 3);

  const batches = asArray(safeJson(batchesRes));
  if (batches.length) {
    const b = pick(batches);
    get(user, `/batch/${b.id}/courses`, "GET /batch/{batchId}/courses");
    get(user, `/attendance/admin/batch/${b.id}`, "GET /attendance/admin/batch/{batchId}");
    get(user, `/assignment/a/batch/${b.id}/submissions?page=1`, "GET /assignment/a/batch/{batchId}/submissions");
  }
  think(1, 3);
}
