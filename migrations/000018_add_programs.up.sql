-- Landing-page programs, fully admin-managed and standalone (no FKs to other tables).
CREATE TABLE programs (
    id          UUID PRIMARY KEY,
    slug        TEXT NOT NULL UNIQUE,
    title       TEXT NOT NULL,
    short_desc  TEXT NOT NULL,
    full_desc   TEXT NOT NULL,
    learn       JSONB NOT NULL DEFAULT '[]'::jsonb, -- array of bullet strings
    age         TEXT NOT NULL DEFAULT '',
    fee         TEXT NOT NULL DEFAULT '',
    level       TEXT NOT NULL DEFAULT '',
    bonus       TEXT NOT NULL DEFAULT '',           -- optional "Included:" line
    accent      TEXT NOT NULL DEFAULT '',           -- optional emoji
    image_url   TEXT NOT NULL DEFAULT '',           -- stable public URL (uploaded) or /programs/*.jpg
    order_no    INT NOT NULL,
    published   BOOLEAN NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_programs_order ON programs(order_no);

INSERT INTO programs (id, slug, title, short_desc, full_desc, learn, age, fee, level, bonus, accent, image_url, order_no, published) VALUES
(gen_random_uuid(), 'junior-stem-explorer', 'Junior STEM Explorer',
 'A fun intro to coding, mental math, and arts & crafts for young learners.',
 'Junior STEM Explorer is designed for children taking their very first steps into technology and creative learning. Through playful activities, hands-on projects, and beginner-friendly tools, students build confidence while developing logical thinking and creativity. This course combines coding fundamentals, mental math exercises, and arts & crafts activities to create a balanced learning experience that keeps young students engaged and excited.',
 '["Beginner coding concepts","ScratchJr activities","Mental math foundations","Creative arts & crafts","Problem-solving skills","Logical thinking"]'::jsonb,
 '7–9 Years', 'SAR 700', 'Beginner', '', '🎨', '/programs/junior-stem-explorer.jpg', 0, TRUE),

(gen_random_uuid(), 'stem-explorer', 'STEM Explorer',
 'Develop coding, logic, teamwork, and problem-solving through exciting STEM.',
 'STEM Explorer helps students strengthen analytical thinking and creativity through coding and problem-solving challenges. Students work on collaborative projects while improving logical reasoning and computational thinking. The course introduces real-world STEM concepts in a simple and engaging way, helping students become confident creators and innovators.',
 '["Coding fundamentals","Problem-solving techniques","Logic building","Team collaboration","Mental math","Creative technology activities"]'::jsonb,
 '10–13 Years', 'SAR 800', 'Intermediate Beginner', '', '🧩', '/programs/stem-explorer.jpg', 1, TRUE),

(gen_random_uuid(), 'digital-creator', 'Digital Creator',
 'Web dev, graphic design, app development, and data analysis through real projects.',
 'Digital Creator is built for students who want to combine creativity with technology. Students can explore multiple digital skills including website development, app design, graphic design, and beginner data analysis. The course encourages project-based learning where students create portfolios, digital artwork, websites, and interactive applications while learning industry-relevant tools.',
 '["Website development","App development","Graphic design","Data analysis basics","Digital creativity","UI/UX fundamentals"]'::jsonb,
 '14–18 Years', 'SAR 800', 'Intermediate to Advanced', 'Includes IGCSE Algebra Support', '💻', '/programs/digital-creator.jpg', 2, TRUE),

(gen_random_uuid(), 'robotics', 'Robotics',
 'Build and program real robots while learning electronics and engineering.',
 'Robotics introduces students to the exciting world of programming, electronics, and engineering through practical robotics projects. Students work with real robotic kits and learn how software and hardware interact together. The course focuses on hands-on experimentation, teamwork, and real-world problem-solving while helping students build confidence in STEM and engineering concepts.',
 '["Robotics programming","Electronics basics","Robot assembly","Problem-solving","Engineering thinking","Hands-on robotics testing"]'::jsonb,
 '11–18 Years', 'SAR 800', 'Beginner to Advanced', 'Robotics kit included with the course', '🤖', '/programs/robotics.jpg', 3, TRUE);
