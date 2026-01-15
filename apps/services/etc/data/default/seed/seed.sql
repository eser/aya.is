-- profiles

insert into "profile" ("created_at", "custom_domain", "deleted_at", "id", "kind", "profile_picture_uri", "pronouns", "properties", "slug", "updated_at", "approved_at") values ('2024-07-07 21:41:56.248504+00', 'eser.dev', NULL, '01J27JA3WM81T7D1TCNJ6Q06P5', 'individual', 'https://avatars.githubusercontent.com/u/866558?v=4', 'he/him', NULL, 'eser', NULL, '2025-07-18 23:41:35.622315+00');
insert into "profile" ("created_at", "custom_domain", "deleted_at", "id", "kind", "profile_picture_uri", "pronouns", "properties", "slug", "updated_at", "approved_at") values ('2024-07-07 17:41:23.73066+00', NULL, NULL, '01J274G56V6VCP6MHK6BFV3AX9', 'product', 'https://objects.aya.is/profile_picture_20250628021500_eserlive.png', NULL, NULL, 'eserlive', NULL, '2025-07-18 23:41:35.622315+00');
insert into "profile" ("created_at", "custom_domain", "deleted_at", "id", "kind", "profile_picture_uri", "pronouns", "properties", "slug", "updated_at", "approved_at") values ('2024-07-07 21:41:56.248504+00', NULL, NULL, '01JT4FS0K93K4GG4Z40ZDFY3EQ', 'organization', 'https://objects.aya.is/profile_picture_20250628021500_aya.png', NULL, NULL, 'aya', NULL, '2025-07-18 23:41:35.622315+00');
insert into "profile" ("created_at", "custom_domain", "deleted_at", "id", "kind", "profile_picture_uri", "pronouns", "properties", "slug", "updated_at", "approved_at") values ('2025-06-27 23:22:24.483038+00', 'yazilimagi.org', NULL, '01JYSV4TD6FSM3F1RKFD4DT5F0', 'organization', 'https://objects.aya.is/profile_picture_20250628021500_yazilimagi.png', NULL, NULL, 'yazilimagi', NULL, '2025-07-18 23:41:35.622315+00');

insert into "profile_tx" ("description", "locale_code", "profile_id", "properties", "title") values ('Gönüllü geliştirilen yazılımlarla oluşan bir yazılım vakfı', 'tr          ', '01JT4FS0K93K4GG4Z40ZDFY3EQ', NULL, 'AYA');
insert into "profile_tx" ("description", "locale_code", "profile_id", "properties", "title") values ('A software foundation formed by volunteer-developed software', 'en          ', '01JT4FS0K93K4GG4Z40ZDFY3EQ', NULL, 'AYA');
insert into "profile_tx" ("description", "locale_code", "profile_id", "properties", "title") values ('Founder @acikyazilim • Streaming @ eser.live • Open Source, DevRel, DevOps and Agile Evangelist • Generalist', 'en          ', '01J27JA3WM81T7D1TCNJ6Q06P5', NULL, 'Eser Ozvataf');
insert into "profile_tx" ("description", "locale_code", "profile_id", "properties", "title") values ('Eser Ozvataf''s YouTube channel', 'en          ', '01J274G56V6VCP6MHK6BFV3AX9', NULL, 'eser.live');
insert into "profile_tx" ("description", "locale_code", "profile_id", "properties", "title") values ('Eser Özvataf''ın YouTube kanalı', 'tr          ', '01J274G56V6VCP6MHK6BFV3AX9', NULL, 'eser.live');
insert into "profile_tx" ("description", "locale_code", "profile_id", "properties", "title") values ('', 'en          ', '01JYSV4TD6FSM3F1RKFD4DT5F0', NULL, 'Yazilim Agi');
insert into "profile_tx" ("description", "locale_code", "profile_id", "properties", "title") values ('Founder @acikyazilim • Streaming @ eser.live • Açık Kaynak, DevRel, DevOps and Agile Evangelist • Generalist', 'tr          ', '01J27JA3WM81T7D1TCNJ6Q06P5', NULL, 'Eser Özvataf');
insert into "profile_tx" ("description", "locale_code", "profile_id", "properties", "title") values ('Network amaçlı çeşitli etkinlikler düzenlemek; topluluğu canlı tutmak amacıyla bir aradayız.', 'tr          ', '01JYSV4TD6FSM3F1RKFD4DT5F0', NULL, 'Yazılım Ağı');

-- profile memberships

insert into "profile_membership" ("created_at", "deleted_at", "finished_at", "id", "kind", "member_profile_id", "profile_id", "properties", "started_at", "updated_at") values ('2024-07-07 17:42:21.029785+00', NULL, NULL, '01J274K5X0E1WYTW3TK80K8WN5', 'owner', '01J27JA3WM81T7D1TCNJ6Q06P5', '01J274G56V6VCP6MHK6BFV3AX9', '{"videos":10}', NULL, NULL);
insert into "profile_membership" ("created_at", "deleted_at", "finished_at", "id", "kind", "member_profile_id", "profile_id", "properties", "started_at", "updated_at") values ('2025-05-11 12:39:55.135571+00', NULL, NULL, '01JTZNQN04K0A4E8WDEDCSNBT4', 'lead', '01J27JA3WM81T7D1TCNJ6Q06P5', '01JT4FS0K93K4GG4Z40ZDFY3EQ', '{"github":{"prs":{"total":45,"resolved":42},"stars":206,"issues":{"total":20,"resolved":10},"commits":156}}', NULL, NULL);
insert into "profile_membership" ("created_at", "deleted_at", "finished_at", "id", "kind", "member_profile_id", "profile_id", "properties", "started_at", "updated_at") values ('2025-06-27 23:38:00.603823+00', NULL, NULL, '01JYSW2EQR9BPQSFWHXJS68XGA', 'owner', '01J27JA3WM81T7D1TCNJ6Q06P5', '01JYSV4TD6FSM3F1RKFD4DT5F0', NULL, NULL, NULL);

-- profile links

insert into "profile_link" ("auth_access_token", "auth_access_token_expires_at", "auth_access_token_scope", "auth_provider", "auth_refresh_token", "auth_refresh_token_expires_at", "created_at", "deleted_at", "id", "is_hidden", "is_managed", "is_verified", "kind", "order", "profile_id", "properties", "public_id", "remote_id", "title", "updated_at", "uri") values (NULL, NULL, NULL, NULL, NULL, NULL, '2025-05-18 21:41:39.087575+00', NULL, '01JVJNGHN8P360BRQMH0T4S6AJ', false, false, true, 'linkedin', 2, '01J27JA3WM81T7D1TCNJ6Q06P5', NULL, NULL, NULL, 'LinkedIn', NULL, 'https://linkedin.com/in/eserozvataf');
insert into "profile_link" ("auth_access_token", "auth_access_token_expires_at", "auth_access_token_scope", "auth_provider", "auth_refresh_token", "auth_refresh_token_expires_at", "created_at", "deleted_at", "id", "is_hidden", "is_managed", "is_verified", "kind", "order", "profile_id", "properties", "public_id", "remote_id", "title", "updated_at", "uri") values (NULL, NULL, NULL, NULL, NULL, NULL, '2025-05-18 21:41:39.087575+00', NULL, '01JVJNGNF0FNW395SX84JEX22F', false, false, true, 'github', 1, '01J27JA3WM81T7D1TCNJ6Q06P5', NULL, NULL, NULL, 'GitHub', NULL, 'https://github.com/eser');
insert into "profile_link" ("auth_access_token", "auth_access_token_expires_at", "auth_access_token_scope", "auth_provider", "auth_refresh_token", "auth_refresh_token_expires_at", "created_at", "deleted_at", "id", "is_hidden", "is_managed", "is_verified", "kind", "order", "profile_id", "properties", "public_id", "remote_id", "title", "updated_at", "uri") values (NULL, NULL, NULL, NULL, NULL, NULL, '2025-04-07 23:09:45.825904+00', NULL, '01JR97YDK41YRX3D9025A0MY9Q', false, false, true, 'x', 3, '01J27JA3WM81T7D1TCNJ6Q06P5', NULL, NULL, NULL, 'X', NULL, 'https://x.com/eserozvataf');
insert into "profile_link" ("auth_access_token", "auth_access_token_expires_at", "auth_access_token_scope", "auth_provider", "auth_refresh_token", "auth_refresh_token_expires_at", "created_at", "deleted_at", "id", "is_hidden", "is_managed", "is_verified", "kind", "order", "profile_id", "properties", "public_id", "remote_id", "title", "updated_at", "uri") values (NULL, NULL, NULL, NULL, NULL, NULL, '2025-05-18 21:29:22.165933+00', NULL, '01JVJMNR1D437XQF6V9TPQ3808', false, false, true, 'youtube', 0, '01J27JA3WM81T7D1TCNJ6Q06P5', NULL, NULL, NULL, 'YouTube', NULL, 'https://youtube.com/@eserlive');
insert into "profile_link" ("auth_access_token", "auth_access_token_expires_at", "auth_access_token_scope", "auth_provider", "auth_refresh_token", "auth_refresh_token_expires_at", "created_at", "deleted_at", "id", "is_hidden", "is_managed", "is_verified", "kind", "order", "profile_id", "properties", "public_id", "remote_id", "title", "updated_at", "uri") values (NULL, NULL, NULL, NULL, NULL, NULL, '2025-06-27 23:49:15.891311+00', NULL, '01JYSWNAHE9CDAGEMN8GYZS158', false, false, true, 'discord', 2, '01JYSV4TD6FSM3F1RKFD4DT5F0', NULL, NULL, NULL, 'Discord', NULL, 'https://discord.com/invite/ckS4huSvEk');
insert into "profile_link" ("auth_access_token", "auth_access_token_expires_at", "auth_access_token_scope", "auth_provider", "auth_refresh_token", "auth_refresh_token_expires_at", "created_at", "deleted_at", "id", "is_hidden", "is_managed", "is_verified", "kind", "order", "profile_id", "properties", "public_id", "remote_id", "title", "updated_at", "uri") values (NULL, NULL, NULL, NULL, NULL, NULL, '2025-06-27 23:47:32.304212+00', NULL, '01JYSWKXKWWAH56YR6G58WBXSV', false, false, true, 'x', 0, '01JYSV4TD6FSM3F1RKFD4DT5F0', NULL, NULL, NULL, 'X', NULL, 'https://x.com/yazilimagiorg');
insert into "profile_link" ("auth_access_token", "auth_access_token_expires_at", "auth_access_token_scope", "auth_provider", "auth_refresh_token", "auth_refresh_token_expires_at", "created_at", "deleted_at", "id", "is_hidden", "is_managed", "is_verified", "kind", "order", "profile_id", "properties", "public_id", "remote_id", "title", "updated_at", "uri") values (NULL, NULL, NULL, NULL, NULL, NULL, '2025-06-27 23:50:24.270027+00', NULL, '01JYSWQHD3FZYTPK2K2AHV3YWG', false, false, true, 'telegram', 1, '01JYSV4TD6FSM3F1RKFD4DT5F0', NULL, NULL, NULL, 'Telegram', NULL, 'https://t.me/yazilimagiorg');

-- profile pages

insert into "profile_page" ("cover_picture_uri", "created_at", "deleted_at", "id", "order", "profile_id", "published_at", "slug", "updated_at") values (NULL, '2025-05-09 20:00:19.402451+00', NULL, '01JT63GNX3TT2RBVHS4D48XJNT', 0, '01JT4FS0K93K4GG4Z40ZDFY3EQ', '2024-10-17 23:27:04.036+00', 'processes', NULL);
insert into "profile_page" ("cover_picture_uri", "created_at", "deleted_at", "id", "order", "profile_id", "published_at", "slug", "updated_at") values (NULL, '2025-05-09 20:03:05.715342+00', NULL, '01JT4GBETQSQ8V7F1RF8RQ987D', 0, '01JT4FS0K93K4GG4Z40ZDFY3EQ', '2024-10-17 23:27:04.036+00', 'policies', NULL);
insert into "profile_page" ("cover_picture_uri", "created_at", "deleted_at", "id", "order", "profile_id", "published_at", "slug", "updated_at") values (NULL, '2025-05-09 20:03:05.715342+00', NULL, '01JT9JHFH5W605XCZ6X4EVZDPC', 0, '01JT4FS0K93K4GG4Z40ZDFY3EQ', '2024-10-17 23:27:04.036+00', 'press', NULL);
insert into "profile_page" ("cover_picture_uri", "created_at", "deleted_at", "id", "order", "profile_id", "published_at", "slug", "updated_at") values (NULL, '2025-05-09 20:03:05.715342+00', NULL, '01JT4M4WJE1FFCC61MF9EZC9XH', 0, '01JT4FS0K93K4GG4Z40ZDFY3EQ', '2024-10-17 23:27:04.036+00', 'about', NULL);
insert into "profile_page" ("cover_picture_uri", "created_at", "deleted_at", "id", "order", "profile_id", "published_at", "slug", "updated_at") values (NULL, '2025-05-09 20:12:18.808159+00', NULL, '01JTRXDZB1Y2AAAK35YJ0S7RXT', 0, '01J27JA3WM81T7D1TCNJ6Q06P5', '2024-10-17 23:27:04.036+00', 'cv', NULL);
insert into "profile_page" ("cover_picture_uri", "created_at", "deleted_at", "id", "order", "profile_id", "published_at", "slug", "updated_at") values (NULL, '2025-06-27 23:41:41.420802+00', NULL, '01JYSW92JKPFHK0JSR1QBX8G5S', 0, '01JYSV4TD6FSM3F1RKFD4DT5F0', '2025-06-27 23:41:41.420802+00', 'channels', NULL);

insert into "profile_page_tx" ("content", "locale_code", "profile_page_id", "summary", "title") values ('## Experience

### Head of Engineering, Teknasyon

May 2023 - Present

### CTO, Datapad

Oct 2022 - May 2023

Datapad is a tech startup is to enable a data-driven culture for teams of all
sizes. They’re on the road to build the world’s easiest data solution. Raised
$1m in pre-seed, backed by 500 Startups.

### Engineering Director, Founding DevRel, Getir

Dec 2021 - Oct 2022

Getir is a decacorn tech startup that delivers ultra-fast groceries.

Achievements:

- The first DevRel of the company, being the founder of the discipline
- Led the “Agile Office”, organizational agility professionals team of agile
  coaches and scrum masters
- Completed 3 bootcamps, along with 14 new recruitments
- Attended some external events
- Celebrated May the 4th: Star Wars Day
- Started technical event series named Getir Tech Talks
- Worked on improvements to hiring processes
- Worked on improvements to onboarding processes
- Crafted an observation-based career framework for the entire technology
  department
- Introduced an internal mobility policy for those who have long tenure
- Crafted a maturity matrix to state common goals for squads in perspective of
  tech-data-product concerns
- Held over 100 face-to-face interviews to collect insights

### Director, Açık Yazılım

Nov 2020 - Dec 2021

Açık Yazılım is founded by myself during the COVID-19 pandemic. As a software
consulting company, it outsources know-how and workforce for its customers
contracted with.

Achievements:

- Founded the company and contracted with customers
- Delivered and maintained a contracted work which is a marketplace specialized
  in holistic wellbeing services
- Provided consultancy to a well-known Turkish Unicorn during their transition
  to product-based agile teams

Technologies, Platforms, and Concepts experienced during the period:

- JavaScript, Python
- Amazon Web Services, Google Cloud Provider, Cloudflare, Cloudflare Pages,
  GitHub Actions
- Kubernetes, Docker
- Node.js, Deno, Cypress, Jest, webpack
- React, React Hooks, GraphQL, Django

### Coordinator of Projects, Turkish Open Source Platform

Dec 2020 - May 2021

An establishment under the auspices of the Ministry of Industry and Technology
of the Turkish state, Turkey Open Source Platform aims to increase the number of
software developers Turkey needs, to increase the qualifications of existing
software developers, and to develop the open-source ecosystem. In this respect,
the platform performs a similar function to CNCF (CloudNative Foundation) and
Linux Foundation.

Achievements:

- Established and facilitated weekly meetings of project groups under five
  different project topics in order to produce the software needed by the public
  and private sectors in Turkey with open source software (These topics were
  Natural Language Processing, Cloud and Database Systems, Mobility, Linux
  (Pardus) and Digital Badges.)
- Formed collaborations between platform stakeholders
- Planned and kicked off online a webinar series and two hackathons
- Designed an architecture and developed project analysis of a digital badge and
  verification project
- Assisted and reviewed technical output related with projects

### Software Manager, Setur

Sep 2018 - Nov 2020

Setur is the tourism and duty-free merchandising company of Koç Holding Group of
Companies, which is on the Fortune 500 list.

Achievements:

- Accomplished the leadership and coaching of Agile & DevOps transformation
  which led the company to have 80% accurate estimations
- Led 9 cross-functional product teams, with a total of 30+ people with various
  functions (UX, QA, SE, BA, etc.)
- Designed company-wide prioritization flow to maintain both product and
  technical roadmaps of products
- Achieved faster onboarding process and development of culture by initiating a
  knowledge base that consists of both cultural and technical guidelines
- Achieved faster troubleshooting and bug detection by established software
  quality assurance and delivery metrics with monitoring, tracing, and alerting
  mechanisms
- Transformed Windows & .NET only ecosystem into polyglot-enabled,
  Linux-containerized, orchestrated and multi-regional infrastructure runs on
  Kubernetes
- Improved IT recruitment processes, implemented technical assessments, attended
  300+ recruitment meetings
- Enabled remote and multi-regional development culture
- Organized a meetup and several training sessions

Technologies, Platforms, and Concepts experienced during the period:

- JavaScript, .NET, Java
- Microsoft Azure, Azure DevOps
- Linux, Windows, macOS
- Kubernetes, Docker
- Postgres, SQL Server, MongoDB, ElasticSearch
- Node.js, Cypress, Jest, webpack, Kibana, Grafana
- React, React Native
- Microservices, DDD
- CloudNative, OpenSource

### Lead Software Developer, KoçSistem

May 2016 - Sep 2018

KoçSistem is the information technology company of Koç Holding Group of
Companies, which is on the Fortune 500 list. They offers a wide range of
products and services in areas such as Internet of Things, Big Data & Analytics,
Security, Corporate Cloud, Corporate Mobility, Business Solutions, Digital
Workforce and RPA.

Achievements:

- Kicked off the company’s main application framework and was the first member
  of the company’s DevOps and Software Architecture teams
- Developed a React project for initializing codebases that use the company’s
  main application framework with selected components on a web wizard form
- Helped other teams to embrace Agile and DevOps practices
- Took parts of analysis and decisions of company’s product choices (such as
  CMSes, orchestration tools, etc)
- Delivered the company’s first serverless project that runs on AWS Lambda
- Delivered and maintained 3 enterprise solutions to customers
- Attended recruitment meetings as a technical interviewer

Technologies, Platforms, and Concepts experienced during the period:

- JavaScript, .NET
- Amazon Web Services
- Jenkins, JMeter, Docker
- SQL Server, Oracle, MongoDB, ElasticSearch
- Node.js, AWS Lambda, AWS Cognito, AWS S3
- React, Redux
- FaaS, DDD

### Software Architect, Mobilex

Sep 2015 - Apr 2016

MobileX is a fast-growing tech company focused on building and marketing highly
scalable mobile applications. Despite being such a young company, they have a
track record of building & growing mobile apps on various verticals that reached
over 120 million users across 160 countries.

Achievements:

- Developed a complete backend and cloud architecture for Instagram-like
  application with location-awareness

Technologies, Platforms, and Concepts experienced during the period:

- JavaScript
- Amazon Web Services
- Terraform, Docker
- Postgres, MongoDB
- Node.js, AWS VPC, AWS EC2, AWS ElasticBeanstalk, AWS S3
- OpenStreetMaps

### E-Commerce Software Consultant, Zaimoglu Holding

Sep 2014 - Aug 2015

### Social Media Oriented Web Development and Research Assistant, Eastern Mediterranean University

Sep 2013 - Sep 2014

### Team Lead, Goldsoft Trading

Jun 2012 - Aug 2013

### Head of Software Development, CMFNet

Jun 2008 - May 2012

### Software Design Engineer, DevBiz

Apr 2006 - Jan 2007

### Backend Developer and System Administrator, Ispro

Mar 2005 - Apr 2006

### Web Developer, Egebilgi

Jun 2004 - Feb 2005

### Web Developer, Sis Hosting

Jan 2001 - Aug 2002

## Education

### Master’s Degree - Information and Communication Technologies in Education

2013 - 2019

Eastern Mediterranean University, CGPA: 3.44 / 4.00

### Bachelor of Science - Information Technology

2010 - 2012

Eastern Mediterranean University, CGPA: 3.12 / 4.00

### Associate’s Degree - Computer Programming

2002 - 2004

Celal Bayar University, CGPA: 78.77/100

## Certificates and Exams

### LEAD Execution Program - Personal Branding, Reputation and Perception Management Training for Leaders

Sep 2019

Issuer: Koç Holding HR - Koç University

### LEAD Execution Program - Customer Experience

Sep 2019

Issuer: Koç Holding HR - Koç University

### Assessment Center

Aug 2019

Issuer: TAG

### Applied Entrepreneurship Training

Apr 2016

Issuer: KOSGEB

### edx CS-169.1x: Software as a Service

Aug 2013

Issuer: Berkeley (University of California)

[https://verify.edx.org/cert/da723362e8b74500992fc4bc84e5324c](https://verify.edx.org/cert/da723362e8b74500992fc4bc84e5324c)

### edx 6.00x: Introduction to Computer Science and Programming

Jun 2013

Issuer: The Massachusetts Institute of Technology

[https://verify.edx.org/cert/4b90bae7e00a4d94b9793c93339334ba](https://verify.edx.org/cert/4b90bae7e00a4d94b9793c93339334ba)

### Academic Personnel and Graduate Education Exam

May 2012

Issuer: OSYM

### General English (EFL) - Upper-Intermediate (Council of Europe B2)

Jun 2011

Issuer: UKLA

### Project Management Methodologies and Practices

Jun 2011

Issuer: Turkcell Akademi

### MCTS - Microsoft Windows Mobile Application Development

May 2010

Issuer: Microsoft

### MCTS - Microsoft Windows Mobile Designing, Implementing and Managing

May 2010

Issuer: Microsoft

[Go back to the homepage](../../README.md)', 'tr          ', '01JTRXDZB1Y2AAAK35YJ0S7RXT', 'Eser Özvataf''ın özgeçmişi', 'Özgeçmiş');
insert into "profile_page_tx" ("content", "locale_code", "profile_page_id", "summary", "title") values ('## Experience

### Head of Engineering, Teknasyon

May 2023 - Present

### CTO, Datapad

Oct 2022 - May 2023

Datapad is a tech startup is to enable a data-driven culture for teams of all
sizes. They’re on the road to build the world’s easiest data solution. Raised
$1m in pre-seed, backed by 500 Startups.

### Engineering Director, Founding DevRel, Getir

Dec 2021 - Oct 2022

Getir is a decacorn tech startup that delivers ultra-fast groceries.

Achievements:

- The first DevRel of the company, being the founder of the discipline
- Led the “Agile Office”, organizational agility professionals team of agile
  coaches and scrum masters
- Completed 3 bootcamps, along with 14 new recruitments
- Attended some external events
- Celebrated May the 4th: Star Wars Day
- Started technical event series named Getir Tech Talks
- Worked on improvements to hiring processes
- Worked on improvements to onboarding processes
- Crafted an observation-based career framework for the entire technology
  department
- Introduced an internal mobility policy for those who have long tenure
- Crafted a maturity matrix to state common goals for squads in perspective of
  tech-data-product concerns
- Held over 100 face-to-face interviews to collect insights

### Director, Açık Yazılım

Nov 2020 - Dec 2021

Açık Yazılım is founded by myself during the COVID-19 pandemic. As a software
consulting company, it outsources know-how and workforce for its customers
contracted with.

Achievements:

- Founded the company and contracted with customers
- Delivered and maintained a contracted work which is a marketplace specialized
  in holistic wellbeing services
- Provided consultancy to a well-known Turkish Unicorn during their transition
  to product-based agile teams

Technologies, Platforms, and Concepts experienced during the period:

- JavaScript, Python
- Amazon Web Services, Google Cloud Provider, Cloudflare, Cloudflare Pages,
  GitHub Actions
- Kubernetes, Docker
- Node.js, Deno, Cypress, Jest, webpack
- React, React Hooks, GraphQL, Django

### Coordinator of Projects, Turkish Open Source Platform

Dec 2020 - May 2021

An establishment under the auspices of the Ministry of Industry and Technology
of the Turkish state, Turkey Open Source Platform aims to increase the number of
software developers Turkey needs, to increase the qualifications of existing
software developers, and to develop the open-source ecosystem. In this respect,
the platform performs a similar function to CNCF (CloudNative Foundation) and
Linux Foundation.

Achievements:

- Established and facilitated weekly meetings of project groups under five
  different project topics in order to produce the software needed by the public
  and private sectors in Turkey with open source software (These topics were
  Natural Language Processing, Cloud and Database Systems, Mobility, Linux
  (Pardus) and Digital Badges.)
- Formed collaborations between platform stakeholders
- Planned and kicked off online a webinar series and two hackathons
- Designed an architecture and developed project analysis of a digital badge and
  verification project
- Assisted and reviewed technical output related with projects

### Software Manager, Setur

Sep 2018 - Nov 2020

Setur is the tourism and duty-free merchandising company of Koç Holding Group of
Companies, which is on the Fortune 500 list.

Achievements:

- Accomplished the leadership and coaching of Agile & DevOps transformation
  which led the company to have 80% accurate estimations
- Led 9 cross-functional product teams, with a total of 30+ people with various
  functions (UX, QA, SE, BA, etc.)
- Designed company-wide prioritization flow to maintain both product and
  technical roadmaps of products
- Achieved faster onboarding process and development of culture by initiating a
  knowledge base that consists of both cultural and technical guidelines
- Achieved faster troubleshooting and bug detection by established software
  quality assurance and delivery metrics with monitoring, tracing, and alerting
  mechanisms
- Transformed Windows & .NET only ecosystem into polyglot-enabled,
  Linux-containerized, orchestrated and multi-regional infrastructure runs on
  Kubernetes
- Improved IT recruitment processes, implemented technical assessments, attended
  300+ recruitment meetings
- Enabled remote and multi-regional development culture
- Organized a meetup and several training sessions

Technologies, Platforms, and Concepts experienced during the period:

- JavaScript, .NET, Java
- Microsoft Azure, Azure DevOps
- Linux, Windows, macOS
- Kubernetes, Docker
- Postgres, SQL Server, MongoDB, ElasticSearch
- Node.js, Cypress, Jest, webpack, Kibana, Grafana
- React, React Native
- Microservices, DDD
- CloudNative, OpenSource

### Lead Software Developer, KoçSistem

May 2016 - Sep 2018

KoçSistem is the information technology company of Koç Holding Group of
Companies, which is on the Fortune 500 list. They offers a wide range of
products and services in areas such as Internet of Things, Big Data & Analytics,
Security, Corporate Cloud, Corporate Mobility, Business Solutions, Digital
Workforce and RPA.

Achievements:

- Kicked off the company’s main application framework and was the first member
  of the company’s DevOps and Software Architecture teams
- Developed a React project for initializing codebases that use the company’s
  main application framework with selected components on a web wizard form
- Helped other teams to embrace Agile and DevOps practices
- Took parts of analysis and decisions of company’s product choices (such as
  CMSes, orchestration tools, etc)
- Delivered the company’s first serverless project that runs on AWS Lambda
- Delivered and maintained 3 enterprise solutions to customers
- Attended recruitment meetings as a technical interviewer

Technologies, Platforms, and Concepts experienced during the period:

- JavaScript, .NET
- Amazon Web Services
- Jenkins, JMeter, Docker
- SQL Server, Oracle, MongoDB, ElasticSearch
- Node.js, AWS Lambda, AWS Cognito, AWS S3
- React, Redux
- FaaS, DDD

### Software Architect, Mobilex

Sep 2015 - Apr 2016

MobileX is a fast-growing tech company focused on building and marketing highly
scalable mobile applications. Despite being such a young company, they have a
track record of building & growing mobile apps on various verticals that reached
over 120 million users across 160 countries.

Achievements:

- Developed a complete backend and cloud architecture for Instagram-like
  application with location-awareness

Technologies, Platforms, and Concepts experienced during the period:

- JavaScript
- Amazon Web Services
- Terraform, Docker
- Postgres, MongoDB
- Node.js, AWS VPC, AWS EC2, AWS ElasticBeanstalk, AWS S3
- OpenStreetMaps

### E-Commerce Software Consultant, Zaimoglu Holding

Sep 2014 - Aug 2015

### Social Media Oriented Web Development and Research Assistant, Eastern Mediterranean University

Sep 2013 - Sep 2014

### Team Lead, Goldsoft Trading

Jun 2012 - Aug 2013

### Head of Software Development, CMFNet

Jun 2008 - May 2012

### Software Design Engineer, DevBiz

Apr 2006 - Jan 2007

### Backend Developer and System Administrator, Ispro

Mar 2005 - Apr 2006

### Web Developer, Egebilgi

Jun 2004 - Feb 2005

### Web Developer, Sis Hosting

Jan 2001 - Aug 2002

## Education

### Master’s Degree - Information and Communication Technologies in Education

2013 - 2019

Eastern Mediterranean University, CGPA: 3.44 / 4.00

### Bachelor of Science - Information Technology

2010 - 2012

Eastern Mediterranean University, CGPA: 3.12 / 4.00

### Associate’s Degree - Computer Programming

2002 - 2004

Celal Bayar University, CGPA: 78.77/100

## Certificates and Exams

### LEAD Execution Program - Personal Branding, Reputation and Perception Management Training for Leaders

Sep 2019

Issuer: Koç Holding HR - Koç University

### LEAD Execution Program - Customer Experience

Sep 2019

Issuer: Koç Holding HR - Koç University

### Assessment Center

Aug 2019

Issuer: TAG

### Applied Entrepreneurship Training

Apr 2016

Issuer: KOSGEB

### edx CS-169.1x: Software as a Service

Aug 2013

Issuer: Berkeley (University of California)

[https://verify.edx.org/cert/da723362e8b74500992fc4bc84e5324c](https://verify.edx.org/cert/da723362e8b74500992fc4bc84e5324c)

### edx 6.00x: Introduction to Computer Science and Programming

Jun 2013

Issuer: The Massachusetts Institute of Technology

[https://verify.edx.org/cert/4b90bae7e00a4d94b9793c93339334ba](https://verify.edx.org/cert/4b90bae7e00a4d94b9793c93339334ba)

### Academic Personnel and Graduate Education Exam

May 2012

Issuer: OSYM

### General English (EFL) - Upper-Intermediate (Council of Europe B2)

Jun 2011

Issuer: UKLA

### Project Management Methodologies and Practices

Jun 2011

Issuer: Turkcell Akademi

### MCTS - Microsoft Windows Mobile Application Development

May 2010

Issuer: Microsoft

### MCTS - Microsoft Windows Mobile Designing, Implementing and Managing

May 2010

Issuer: Microsoft

[Go back to the homepage](../../README.md)', 'en          ', '01JTRXDZB1Y2AAAK35YJ0S7RXT', 'CV of Eser Ozvataf', 'Curriculum Vitae');
insert into "profile_page_tx" ("content", "locale_code", "profile_page_id", "summary", "title") values ('**AYA** (Açık Yazılım Ağı) is a platform enabling volunteers motivated by social good to collaborate and develop quality information technology solutions.

We have collectively realized projects by providing expert volunteer support and technical resources. These projects have been open-source solutions accessible and controllable by everyone, primarily official institutions and non-governmental organizations. Our role has been facilitating volunteer communication and resource management coordination.

As a non-profit organization focusing on social benefit, AYA continues its journey with a [software foundation identity](/aya/stories/20240107-vakif-kimligimiz-ve-sinirlarimiz) hosting various IT projects.

<iframe className="w-full aspect-[16/9] mb-8" width="560" height="505" src="https://www.youtube.com/embed/mktym-FxWvg" title="Yazılımcı Günleri 2023 | ESER ÖZVATAF - AÇIK YAZILIM AĞI VE AFET HARİTA PROJESİ" frameBorder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" allowFullScreen></iframe>

## Brief History

Today, software development practices have adopted open-source and free software principles, enabling collaborative IT solutions online.

Platforms like GitHub allow initiating and contributing to open-source projects easily. This development results from years of work by foundations such as Mozilla, Linux Foundation, Apache Foundation, and CNCF.

In Türkiye, several open-source communities have emerged, such as the Linux Users Association and Pardus, contributing significantly. Among these was the [Açık Kaynak İnisiyatifi](https://acikkaynak.github.io), founded in 2015, promoting open-source development through GitHub.

After the February 2023 earthquake in Kahramanmaraş, volunteers coordinated through Twitter, creating a platform that quickly mobilized over 24,000 IT volunteers via GitHub and Discord.

Named **AYA**, this initiative enabled nearly 50 projects addressing critical post-earthquake issues, liaising with organizations like AFAD, AHBAP, and AKUT.

## Objectives

Our primary aim is addressing IT needs of field organizations during crises, combining optimal resources and expertise to minimize chaos and efficiently direct volunteer efforts.

Based on our experiences, we revised our purpose to:
- Raise awareness about IT solutions for social good.
- Lead development of practical IT solutions.
- Contribute to a sustainable open-source ecosystem.
- Promote widespread use of these solutions.

## Goals

Our current goals are:
- Increasing IT volunteer participation.
- Encouraging developer communication.
- Cultivating collective work culture.
- Supporting development of social-benefit IT solutions.
- Meeting institutional needs effectively.
- Applying industry best practices.
- Offering projects to help developers gain experience during industry bottlenecks.', 'en          ', '01JT4M4WJE1FFCC61MF9EZC9XH', 'Information about AYA', 'About');
insert into "profile_page_tx" ("content", "locale_code", "profile_page_id", "summary", "title") values ('**AYA** (Açık Yazılım Ağı), sosyal bir yarar oluşturmak motivasyonuyla nitelikli bilişim çözümleri
üreten gönüllülerin bir araya gelmesine ve projeler gerçekleştirmesine imkan sağlayan bir platformdur.

Bugüne kadar projelere yetkin ve uzman gönüllü desteği ve teknik kaynaklar sunarak kolektif bir
şekilde gerçeklenmelerini sağladık. Bu projeler resmi kurumlar ve sivil toplum kuruluşları başta olmak
üzere herkesin kullanımına ve denetimine açık, açık kaynak lisanslarına sahip bilişim çözümleri
oldular. Bizim buradaki rolümüz hem gönüllülerin aralarındaki iletişime aracı olmak hem de kaynak
yönetiminde koordinasyon sağlamak oldu.

Sosyal yarar gözeten ve kar amacı gütmeyen bir oluşum olarak AYA, bugün yoluna portföyünde bilişim
projeleri barındıran [yazılım vakfı kimliğiyle](/aya/stories/20240107-vakif-kimligimiz-ve-sinirlarimiz)
devam ediyor.

<iframe
    className="w-full aspect-[16/9] mb-8"
    width="560"
    height="505"
    src="https://www.youtube.com/embed/mktym-FxWvg"
    title="Yazılımcı Günleri 2023 | ESER ÖZVATAF - AÇIK YAZILIM AĞI VE AFET HARİTA PROJESİ"
    frameBorder="0"
    allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share"
    allowFullScreen
>
</iframe>


## Kısa Geçmişimiz

Günümüzde yazılım geliştirme pratikleri konjonktürel olarak açık kaynak ve özgür yazılım
anlayışlarını özümseyerek, uzmanların internet üzerinden birlikte bilişim çözümleri
geliştirebilmelerine elverişli bir noktaya gelmiştir.

GitHub gibi mecralar üzerinde açık kaynaklı yazılım projeleri başlatabilmek ve mevcut projelere
rahatlıkla katkılarda bulunabilmek mümkün. Gelinen bu nokta Mozilla, Linux Foundation, Apache
Foundation, Cloud Native Computing Foundation vb. yazılım vakıflarının yıllarca süren
çalışmalarının bir sonucu olarak görülebilir.

Ülkemizde de geçmişte Linux Kullanıcı Derneği, Pardus vb. oluşumlar etrafında birçok açık kaynak
topluluğu oluşmuş, çeşitli çalışmalara öncülük etmişlerdir. Bunlardan biri olan ve 2015’de kurulan
[Açık Kaynak İnisiyatifi](https://acikkaynak.github.io) GitHub üzerinden insanları bir araya getirmeyi planlamış,
açık kaynak geliştirme alanında insanların birbiri ile iletişimde olmasını sağlamayı amaçlamıştır.
Açık Kaynak İnisiyatifinin yürüttüğü bu çalışmalar hem açık kaynak gönüllüsü insanların birbiri ile
iletişımde olması, hem de bilgi birikimi oluşturmaları konusunda bayrağı ileri taşımıştır.

Şubat 2023’te Kahramanmaraş’ta gerçekleşen ve 10 ili etkileyen deprem sonrası, Twitter üzerinde
koordine olan gönüllülerin oluşturduğu bir grup da buradaki bilgi birikimini kullanarak “deprem
yardımı için” organize olmuş; kısa sürede çalışmaları GitHub ve Discord üzerinde koordine ederek,
24.000’den fazla bilişim sektörü gönüllüsünü iletişimde tutan ve kaynak yönetimini üstlenen bir
platform oluşturmuştur.

Bir süre sonra **AYA** olarak isimlendirdiğimiz bu oluşum, topluluk tarafından GitHub
üzerinde 50’ye yakın proje çıktısı oluşturulmasına imkan sağladı; deprem sonrası hayati önem
taşıyan sorunları çözebilecek teknolojiler sunabilmek için AFAD, AHBAP, AKUT gibi birçok resmi
kurum ve sivil toplum kuruluşu ile irtibat noktası olduk ve bilişim desteğinin sağlanmasında aracı
olarak hareket ettik.


## Amacımız

Temel çerçevede kriz ve ihtiyaç anlarında sahada olan ve toplum yararına çalışan kurumların bilişim
ihtiyaçlarını karşılamak için bir araya geldik. Devamında, alandaki en faydalı katkıları
sağlayabilecek insanları ve teknolojik kaynakları bir araya getirmeyi amaçlayarak yola çıktık.

Böylece hem kaynakları verimli kullanmak, afet senaryoları gibi olası kriz anlarındaki kaosu en aza
indirmek, hem de gönüllülerin çabalarını doğru yönlendirmek için bir platform sağlayabiliyorduk.
Ancak [bu süreçte](/aya/stories/20230217-acik-yazilim-agi-ve-6-subat) edindiğimiz deneyim bizlere gösterdi
ki, kriz anları hareket etmek için çok geç kalınmış anlar oluyor.

Bu nedenle amacımızı revize ederek, geçmişte [Açık Kaynak İnisiyatifi](https://acikkaynak.github.io)nin
oluşturduğu kapsamdan devam etmeye karar verdik. Bu kapsamda AYA’nın amacı;

- Sosyal yarar oluşturacak bilişim çözümleri geliştirmeye yönelik farkındalık oluşturmak,
- pratik ihtiyaçların karşılanması için bilgi birikimini ve teknik kaynakları bir araya getirerek
  bilişim çözümlerin tasarlanmasına öncü olmak,
- sürdürülebilir bir “açık kaynak ekosisteminin” oluşmasına ve açık kaynak projeler üretilmesine
  katkıda bulunmak,
- üretilen projelerin sahalarda kullanımı ve yaygınlaşmasını sağlamak,

olarak belirledik.


## Hedeflerimiz

Açık kaynak ve özgür yazılım pratiklerinin sosyal yarar oluşturmak için kullanılmasına yönelik
bir örnek teşkil etmek hedeflerimizin çatısını oluşturuyor.

Mevcut Hedeflerimiz,
- Bilişim alanındaki gönüllü sayısını arttırmak,
- Geliştiricilerin birbirleri ile iletişimini sağlamak,
- Kolektif çalışma kültürünü geliştirmek,
- Sosyal fayda üreten bilişim çözümlerinin üretimini teşvik etmek,
- Toplum yararına çalışan kurumların ihtiyaçlarını karşılamak,
- Sektörde en güncel ve geçerli pratiklerinin üretilen bilişim çözümlerinde uygulanmasını
  sağlamak,
- Yazılım sektöründe darboğaz oluştuğunda geliştiricilerin deneyim kazanmasını sağlayacak
  projeler bulundurmak,

olarak listelenebilir.', 'tr          ', '01JT4M4WJE1FFCC61MF9EZC9XH', 'AYA hakkında bilgi', 'Hakkında');
insert into "profile_page_tx" ("content", "locale_code", "profile_page_id", "summary", "title") values ('<List>
  <ListItem item={{
  "title": "Yazılımcı Günleri 2023",
  "description": "AYA ve Afet Harita Projesi",
  "url": "https://www.youtube.com/watch?v=mktym-FxWvg",
  "host": "FolksDev",
  "participants": ["Eser Özvataf"],
  "date": "2023-08-12"
 }} />
 <ListItem item={{
  "title": "Kesin Yaşanmıştır",
  "description": "Bir günde 20 bin gönüllü topladık",
  "url": "https://www.youtube.com/watch?v=1ahIvlbVVbs",
  "host": "Teknasyon",
  "participants": ["Fatih Kadir Akın"],
  "date": "2023-07-04"
 }} />
 <ListItem item={{
  "title": "Nevşin Mengü YouTube",
  "description": "Seçim güvenliğinde yer almak isteyenler bu videoyu izleyin!",
  "url": "https://www.youtube.com/watch?v=jWJqhQYy7Ow",
  "host": "Nevşin Mengü",
  "participants": [],
  "date": "2023-05-26"
 }} />
 <ListItem item={{
  "title": "Turk-internet.com",
  "description": "Gönüllü Gözlemci Olarak Gideceğiniz Yeri Bu Haritadan Bulun",
  "url": "https://turk-internet.com/gonullu-gozlemci-olarak-gideceginiz-yeri-bu-haritadan-bulun/",
  "host": "Turk-internet.com",
  "participants": ["Füsun Sarp Nebil"],
  "date": "2023-05-26"
 }} />
 <ListItem item={{
  "title": "Açık Veri Zirvesi 2023",
  "description": "Gönüllü Gözlemci Olarak Gideceğiniz Yeri Bu Haritadan Bulun",
  "url": "https://turk-internet.com/gonullu-gozlemci-olarak-gideceginiz-yeri-bu-haritadan-bulun/",
  "host": "Açık Veri ve Teknoloji Derneği",
  "participants": ["Fatih Kadir Akın"],
  "date": "2023-05-25"
 }} />
 <ListItem item={{
  "title": "About Amazon",
  "description": "Amazon has activated its disaster relief services to help organizations responding to the devastating earthquake in Türkiye and Syria.",
  "url": "https://www.aboutamazon.eu/news/community-engagement/how-amazon-is-helping-communities-impacted-by-the-earthquake-in-turkey",
  "host": "Amazon Team",
  "participants": [],
  "date": "2023-02-27"
 }} />
 <ListItem item={{
  "title": "Webtekno röportajımız",
  "description": "23 Bin Gönüllü Yazılımcıdan Hayat Kurtaran Platform! Afet.org",
  "url": "https://www.youtube.com/watch?v=2rmDUQYQc3k",
  "host": "Webtekno",
  "participants": ["Melihcan Örüş"],
  "date": "2023-02-16"
 }} />
 <ListItem item={{
  "title": "Google Türkiye Resmi Blog Sitesi",
  "description": "Yaşanan Deprem Felaketinin Ardından",
  "url": "https://turkiye.googleblog.com/2023/02/yasanan-deprem-felaketinin-ardndan.html",
  "host": "Google",
  "participants": [],
  "date": "2023-02-15"
 }} />
 <ListItem item={{
  "title": "Ekonomim",
  "description": "Girişimci Z kuşağı tarih yazıyor",
  "url": "https://www.ekonomim.com/kose-yazisi/girisimci-z-kusagi-tarih-yaziyor/682924",
  "host": "Ekonomim",
  "participants": ["D. Ferhat Demir"],
  "date": "2023-02-15"
 }} />
 <ListItem item={{
  "title": "Milliyet",
  "description": "Yazılım ordusu seferber oldu!",
  "url": "https://www.milliyet.com.tr/yazarlar/hanife-bas/yazilim-ordusu-seferber-oldu-6902301",
  "host": "Milliyet",
  "participants": ["Hanife Baş"],
  "date": "2023-02-12"
 }} />
 <ListItem item={{
  "title": "TIME",
  "description": "Tech Volunteers Rush to Save Turkey’s Earthquake Survivors",
  "url": "https://time.com/6254500/turkey-earthquake-twitter-musk-rescue/",
  "host": "TIME",
  "participants": ["Vera Bergengruen"],
  "date": "2023-02-11"
 }} />
 <ListItem item={{
  "title": "Euronews",
  "description": "Turkey-Syria earthquakes: How Twitter has helped find survivors trapped beneath the rubble",
  "url": "https://www.euronews.com/next/2023/02/10/how-twitter-helped-find-survivors-trapped-beneath-rubble-after-turkeys-earthquakes",
  "host": "Euronews",
  "participants": ["Aylin Elci"],
  "date": "2023-02-10"
 }} />
 <ListItem item={{
  "title": "WIRED",
  "description": "Tech Volunteers Rush to Save Turkey’s Earthquake Survivors",
  "url": "https://www.wired.com/story/tech-volunteers-rush-to-save-turkeys-earthquake-survivors/",
  "host": "WIRED",
  "participants": ["Robyn Huang"],
  "date": "2023-02-08"
 }} />
</List>', 'tr          ', '01JT9JHFH5W605XCZ6X4EVZDPC', 'Basında hakkımızda yapılan haberler', 'Basında Biz');
insert into "profile_page_tx" ("content", "locale_code", "profile_page_id", "summary", "title") values ('Henüz içerik bulunmamaktadır.', 'tr          ', '01JT4GBETQSQ8V7F1RF8RQ987D', 'AYA kuralları', 'Kurallar');
insert into "profile_page_tx" ("content", "locale_code", "profile_page_id", "summary", "title") values ('Content not yet available.', 'en          ', '01JT4GBETQSQ8V7F1RF8RQ987D', 'AYA policies', 'Policies');
insert into "profile_page_tx" ("content", "locale_code", "profile_page_id", "summary", "title") values ('<List>
  <ListItem item={{
  "title": "Yazılımcı Günleri 2023",
  "description": "AYA ve Afet Harita Projesi",
  "url": "https://www.youtube.com/watch?v=mktym-FxWvg",
  "host": "FolksDev",
  "participants": ["Eser Özvataf"],
  "date": "2023-08-12"
 }} />
 <ListItem item={{
  "title": "Kesin Yaşanmıştır",
  "description": "Bir günde 20 bin gönüllü topladık",
  "url": "https://www.youtube.com/watch?v=1ahIvlbVVbs",
  "host": "Teknasyon",
  "participants": ["Fatih Kadir Akın"],
  "date": "2023-07-04"
 }} />
 <ListItem item={{
  "title": "Nevşin Mengü YouTube",
  "description": "Seçim güvenliğinde yer almak isteyenler bu videoyu izleyin!",
  "url": "https://www.youtube.com/watch?v=jWJqhQYy7Ow",
  "host": "Nevşin Mengü",
  "participants": [],
  "date": "2023-05-26"
 }} />
 <ListItem item={{
  "title": "Turk-internet.com",
  "description": "Gönüllü Gözlemci Olarak Gideceğiniz Yeri Bu Haritadan Bulun",
  "url": "https://turk-internet.com/gonullu-gozlemci-olarak-gideceginiz-yeri-bu-haritadan-bulun/",
  "host": "Turk-internet.com",
  "participants": ["Füsun Sarp Nebil"],
  "date": "2023-05-26"
 }} />
 <ListItem item={{
  "title": "Açık Veri Zirvesi 2023",
  "description": "Gönüllü Gözlemci Olarak Gideceğiniz Yeri Bu Haritadan Bulun",
  "url": "https://turk-internet.com/gonullu-gozlemci-olarak-gideceginiz-yeri-bu-haritadan-bulun/",
  "host": "Açık Veri ve Teknoloji Derneği",
  "participants": ["Fatih Kadir Akın"],
  "date": "2023-05-25"
 }} />
 <ListItem item={{
  "title": "About Amazon",
  "description": "Amazon has activated its disaster relief services to help organizations responding to the devastating earthquake in Türkiye and Syria.",
  "url": "https://www.aboutamazon.eu/news/community-engagement/how-amazon-is-helping-communities-impacted-by-the-earthquake-in-turkey",
  "host": "Amazon Team",
  "participants": [],
  "date": "2023-02-27"
 }} />
 <ListItem item={{
  "title": "Webtekno röportajımız",
  "description": "23 Bin Gönüllü Yazılımcıdan Hayat Kurtaran Platform! Afet.org",
  "url": "https://www.youtube.com/watch?v=2rmDUQYQc3k",
  "host": "Webtekno",
  "participants": ["Melihcan Örüş"],
  "date": "2023-02-16"
 }} />
 <ListItem item={{
  "title": "Google Türkiye Resmi Blog Sitesi",
  "description": "Yaşanan Deprem Felaketinin Ardından",
  "url": "https://turkiye.googleblog.com/2023/02/yasanan-deprem-felaketinin-ardndan.html",
  "host": "Google",
  "participants": [],
  "date": "2023-02-15"
 }} />
 <ListItem item={{
  "title": "Ekonomim",
  "description": "Girişimci Z kuşağı tarih yazıyor",
  "url": "https://www.ekonomim.com/kose-yazisi/girisimci-z-kusagi-tarih-yaziyor/682924",
  "host": "Ekonomim",
  "participants": ["D. Ferhat Demir"],
  "date": "2023-02-15"
 }} />
 <ListItem item={{
  "title": "Milliyet",
  "description": "Yazılım ordusu seferber oldu!",
  "url": "https://www.milliyet.com.tr/yazarlar/hanife-bas/yazilim-ordusu-seferber-oldu-6902301",
  "host": "Milliyet",
  "participants": ["Hanife Baş"],
  "date": "2023-02-12"
 }} />
 <ListItem item={{
  "title": "TIME",
  "description": "Tech Volunteers Rush to Save Turkey’s Earthquake Survivors",
  "url": "https://time.com/6254500/turkey-earthquake-twitter-musk-rescue/",
  "host": "TIME",
  "participants": ["Vera Bergengruen"],
  "date": "2023-02-11"
 }} />
 <ListItem item={{
  "title": "Euronews",
  "description": "Turkey-Syria earthquakes: How Twitter has helped find survivors trapped beneath the rubble",
  "url": "https://www.euronews.com/next/2023/02/10/how-twitter-helped-find-survivors-trapped-beneath-rubble-after-turkeys-earthquakes",
  "host": "Euronews",
  "participants": ["Aylin Elci"],
  "date": "2023-02-10"
 }} />
 <ListItem item={{
  "title": "WIRED",
  "description": "Tech Volunteers Rush to Save Turkey’s Earthquake Survivors",
  "url": "https://www.wired.com/story/tech-volunteers-rush-to-save-turkeys-earthquake-survivors/",
  "host": "WIRED",
  "participants": ["Robyn Huang"],
  "date": "2023-02-08"
 }} />
</List>', 'en          ', '01JT9JHFH5W605XCZ6X4EVZDPC', 'News about us in the press', 'In the Press');
insert into "profile_page_tx" ("content", "locale_code", "profile_page_id", "summary", "title") values ('## General Structure

AYA is a [software foundation](/aya/stories/20240107-vakif-kimligimiz-ve-sinirlarimiz). It hosts independently progressing projects.

## Projects

Projects at AYA are led by appointed project leaders. Each project leader forms a project team at their discretion.

However, the model encouraged by AYA is for previously established communities such as [kamp-us](https://github.com/kamp-us), [patika.dev](https://github.com/patikadev), [eser.live](https://eser.live), etc., to lead these projects.

Project leaders may request specialist resources or technical requirements from AYA based on the project''s needs. These requests are fulfilled from AYA’s resource pools.

New project proposals should be submitted via the [Project Idea Form](https://forms.gle/MJFHKcp2yhTpjeod7).

## Project Team

Each project team has a project leader. Permanent team members who contribute to the project are called maintainers.

However, similar to open-source projects, contributors do not need to be part of the project team at AYA. External contributions are reviewed and decided upon by the fixed members of the project team.

To join project teams, ensure your information is listed in our expert pool via the [Volunteer Form](https://forms.gle/o8SWkeGB5n4SdBnr8).

## Project Development Process

The project development process follows procedures set by the project leader. Each team may have its own autonomous standards. AYA’s only expectation here is adherence to the fundamental criteria we have defined for our projects.

### Our Fundamental Criteria

- Must be published as open-source with a Permissive or CopyLeft license.
- If the project is tied to any domain, ownership of this domain must be transferable.
- An action plan must be prepared to swiftly deploy the project when needed.
- Not storing data, nor violating any legal obligations.
- Adherence to our [Code of Conduct](https://github.com/aya/policies).

### Project Evaluations

Projects are evaluated according to AYA’s schedule availability. Decisions like incubation, providing support, and graduation of projects are made following these evaluations.', 'en          ', '01JT63GNX3TT2RBVHS4D48XJNT', 'How processes work at AYA?', 'Processes');
insert into "profile_page_tx" ("content", "locale_code", "profile_page_id", "summary", "title") values ('## Genel Çatı

AYA bir [yazılım vakfıdır](/aya/stories/20240107-vakif-kimligimiz-ve-sinirlarimiz).
Bünyesinde birbirinden bağımsız ilerleyen projeler bulunur.


## Projeler

AYA’da projeler, projelerin atanan liderleri tarafından yürütülür. Her proje için proje
liderinin tasarruflarıyla bir proje ekibi oluşturulur.

Ancak AYA’nın teşvik ettiği model, [kamp-us](https://github.com/kamp-us),
[patika.dev](https://github.com/patikadev), [eser.live](https://eser.live)
vb. önceden oluşturulmuş toplulukların projeleri yürütüyor olmasıdır.

Proje lideri projenin ihtiyaçlarına istinaden AYA’dan uzman kaynağı veya teknik ihtiyaçlar
talep edebilir. Bu talepler, AYA’nın kaynak havuzlarından karşılanır.

Yeni proje önerilerini [Proje Fikir Formu](https://forms.gle/MJFHKcp2yhTpjeod7) aracılığı
ile iletilmelidir.


## Proje Ekibi

Her proje ekibinde bir proje lideri bulunur. Projenin ekibinde kalıcı olarak yer alan ve
katkıda bulunan kimselere ise proje üyesi (maintainer) denir.

Ancak açık kaynak projelerinde olduğu gibi, AYA’da da proje katkıcısı (contributor) olmak
için proje ekibinde yer almak gerekmez. Dışarıdan gelen katkılar, proje ekibinin sabit
üyeleri tarafından değerlendirilir ve projeye dahil edilmesine karar verilir.

Proje ekiplerinde yer almak için [Gönüllü Formu](https://forms.gle/o8SWkeGB5n4SdBnr8)
aracılığı uzman havuzumuzda bilgilerinizin yer almasını sağlamalısınız.


## Proje Geliştirme Süreci

Proje geliştirme süreci, proje liderinin belirlediği süreçlere göre yürütülür. Her ekibin
kendi otonom standartları olabilir. AYA’nın buradaki tek beklentisi projelerimiz için
belirlediğimiz temel kriterlere uyulmasıdır.

### Temel Kriterlerimiz

- Permissive veya CopyLeft bir lisans ile açık kaynak olarak yayınlanması,
- Projenin herhangi bir alan adına bağlı olması durumunda bu alan adının
  sahipliğinin devredilebilir olması,
- İhtiyaç olduğunda projenin hızla devreye sokulabilmesi için gerekli eylem
  planının hazır olması,
- Veri tutmamak, herhangi bir yasal yükümlülüğü ihlal etmemek,
- [Davranış Kodumuz](https://github.com/aya/policies)a uygun davranmak,

### Proje Değerlendirmeleri

Projelerin değerlendirilmesi AYA’nın kendi zaman uygunluğuna göre yapılır. Değerlendirme
sonucunda projelerin kuluçkaya alınması, destek verilmesi ve mezun edilmesi gibi kararlar
alınır.', 'tr          ', '01JT63GNX3TT2RBVHS4D48XJNT', 'AYA’da süreçler nasıl işliyor?', 'Süreçler');
insert into "profile_page_tx" ("content", "locale_code", "profile_page_id", "summary", "title") values ('Content not yet available.', 'en          ', '01JYSW92JKPFHK0JSR1QBX8G5S', 'Our communication channels', 'Channels');
insert into "profile_page_tx" ("content", "locale_code", "profile_page_id", "summary", "title") values ('## Yazılım Ağı

**Duyuru Kanalları**

- 🐦 [twitter/x](https://x.com/yazilimagiorg): Twitter/X hesabımız
- 📣 [yazılım ağı duyuruları](https://t.me/yazilimagiorg): yazılım ağı genelinde yapılan etkinlik, yayın v.b. konular için telegram duyuru kanalı.
- 🎙️ [eser.live ve yazılım ağı](https://t.me/eserlive): eser.live bünyesinde hazırlanan podcast, canlı yayın v.b. içerikler için duyuru kanalı.

**Tartışma Kanalları**

- 💬 [discord](https://discord.com/invite/ckS4huSvEk): Discord sunucumuz
- 🗣️ [yazılım ağı tartışma (meta-genel kanal)](https://t.me/yazilimagi_tartisma): yazılım ağı genelinde, herkesin yazabileceği genel tartışma kanalımız.
- 🌆 [istanbul yazılım ağı](https://t.me/yazilimagi_istanbul): İstanbul özelindeki yazılım ağı kanalımız.
- 🏖️ [izmir yazılım ağı](https://t.me/yazilimagi_izmir): İzmir özelindeki yazılım ağı kanalımız.

**İlgi Alanı Kanalları**

- ⌨️ [MechIstanbul](https://t.me/mechistanbul): Mekanik Klavye Topluluğu
- 🎮 [Oyuncu Yazılımcılar](https://t.me/oyunyazilim): Oyun Sevenler
- 📚 [Bilim Kurgu ve/veya Fantastik Edebiyat Yazılımcıları](https://t.me/bilimkurguyazilim): Bilim-Kurgu/Fantastik Edebiyat Sevenler
- ☕ [Kahve Demleyen Yazılımcılar](https://t.me/kahveciyazilimcilar): Kahve Sevenler
- 🍷 [İçki Tadan Yazılımcılar](https://t.me/ickiyazilim): Alkol Sevenler
- 📸 [Fotoğraf Çeken Yazılımcılar](https://t.me/fotoyazilim): Fotoğraf Çekimiyle İlgilenenler
- 🎬 [Sinefil ve Dizifil Yazılımcılar](https://t.me/sinefilyazilim): Sinema ve Dizilerle İlgilenenler
', 'tr          ', '01JYSW92JKPFHK0JSR1QBX8G5S', 'İletişim kanallarımız', 'Kanallar');

-- stories

insert into "story" ("author_profile_id", "created_at", "deleted_at", "id", "is_featured", "kind", "properties", "slug", "status", "story_picture_uri", "updated_at") values ('01J27JA3WM81T7D1TCNJ6Q06P5', '2024-09-04 22:24:46.979119+00', NULL, '01J6ZJ7SSNZ0JT71YTCB1YJ8F2', false, 'article', NULL, '20190531-wordpressin-kerameti', 'published', NULL, NULL);
insert into "story" ("author_profile_id", "created_at", "deleted_at", "id", "is_featured", "kind", "properties", "slug", "status", "story_picture_uri", "updated_at") values ('01J27JA3WM81T7D1TCNJ6Q06P5', '2024-09-04 22:24:44.520546+00', NULL, '01J6ZJ7QCTXKKS7DJ9NYMC9ZSH', false, 'article', NULL, '20161211-turkiyedeki-acik-kaynak-toplulugu', 'published', NULL, NULL);
insert into "story" ("author_profile_id", "created_at", "deleted_at", "id", "is_featured", "kind", "properties", "slug", "status", "story_picture_uri", "updated_at") values ('01J27JA3WM81T7D1TCNJ6Q06P5', '2024-09-04 22:24:46.827225+00', NULL, '01J6ZJ7SMXPB77HCZXSAQBN57W', false, 'article', NULL, '20221005-devrel-02', 'published', NULL, NULL);
insert into "story" ("author_profile_id", "created_at", "deleted_at", "id", "is_featured", "kind", "properties", "slug", "status", "story_picture_uri", "updated_at") values ('01J27JA3WM81T7D1TCNJ6Q06P5', '2024-09-04 22:24:42.741007+00', NULL, '01J6ZJ7NN6WHYB587G2TXRTS3Q', false, 'article', NULL, '20111128-cloud-computingde-son-kullanici-acisindan-gizlilik', 'published', NULL, NULL);
insert into "story" ("author_profile_id", "created_at", "deleted_at", "id", "is_featured", "kind", "properties", "slug", "status", "story_picture_uri", "updated_at") values ('01J27JA3WM81T7D1TCNJ6Q06P5', '2024-09-04 22:24:41.899933+00', NULL, '01J6ZJ7MTZNZ207P1NYAA41A61', false, 'article', NULL, '20170226-uretmek-ve-gonulluluk', 'published', NULL, NULL);
insert into "story" ("author_profile_id", "created_at", "deleted_at", "id", "is_featured", "kind", "properties", "slug", "status", "story_picture_uri", "updated_at") values ('01J27JA3WM81T7D1TCNJ6Q06P5', '2024-09-04 22:24:42.877668+00', NULL, '01J6ZJ7NSFX726JYRV59WX72MP', false, 'article', NULL, '20221016-ogrenmek-icin-o-ani-beklemek', 'published', NULL, NULL);
insert into "story" ("author_profile_id", "created_at", "deleted_at", "id", "is_featured", "kind", "properties", "slug", "status", "story_picture_uri", "updated_at") values ('01J27JA3WM81T7D1TCNJ6Q06P5', '2024-09-04 22:24:42.461191+00', NULL, '01J6ZJ7NCGNHWTP081Q3DD0BXN', false, 'article', NULL, '20081231-wpfde-font-rendering', 'published', NULL, NULL);
insert into "story" ("author_profile_id", "created_at", "deleted_at", "id", "is_featured", "kind", "properties", "slug", "status", "story_picture_uri", "updated_at") values ('01J27JA3WM81T7D1TCNJ6Q06P5', '2024-09-04 22:24:42.037763+00', NULL, '01J6ZJ7MZ7DF5QEDZQFM7K8NV5', false, 'article', NULL, '20160221-hayal-kirikligina-ugramamak-icin-sorgulanmasi-gerekenler', 'published', NULL, NULL);
insert into "story" ("author_profile_id", "created_at", "deleted_at", "id", "is_featured", "kind", "properties", "slug", "status", "story_picture_uri", "updated_at") values ('01J27JA3WM81T7D1TCNJ6Q06P5', '2024-09-04 22:24:43.37063+00', NULL, '01J6ZJ7P8W06V3WY2R70S53E8W', false, 'article', NULL, '20160221-90larda-programlama-ogrenmek', 'published', NULL, NULL);
insert into "story" ("author_profile_id", "created_at", "deleted_at", "id", "is_featured", "kind", "properties", "slug", "status", "story_picture_uri", "updated_at") values ('01J27JA3WM81T7D1TCNJ6Q06P5', '2024-09-04 22:24:46.537267+00', NULL, '01J6ZJ7SBVZ2NHN2BE00GT00ZF', false, 'article', NULL, '20161127-yarini-okuyabilir-miyiz', 'published', NULL, NULL);
insert into "story" ("author_profile_id", "created_at", "deleted_at", "id", "is_featured", "kind", "properties", "slug", "status", "story_picture_uri", "updated_at") values ('01J27JA3WM81T7D1TCNJ6Q06P5', '2024-09-04 22:24:45.513668+00', NULL, '01J6ZJ7RBVQ3GD7H1NHN0ZYBSY', false, 'article', NULL, '20161126-visual-studio-for-mac-deneyimlerim', 'published', NULL, NULL);
insert into "story" ("author_profile_id", "created_at", "deleted_at", "id", "is_featured", "kind", "properties", "slug", "status", "story_picture_uri", "updated_at") values ('01J27JA3WM81T7D1TCNJ6Q06P5', '2024-09-04 22:24:44.096282+00', NULL, '01J6ZJ7PZJWHBFYY8Y5A0FSM49', false, 'article', NULL, '20230724-ogrenme-yol-haritasi', 'published', NULL, NULL);
insert into "story" ("author_profile_id", "created_at", "deleted_at", "id", "is_featured", "kind", "properties", "slug", "status", "story_picture_uri", "updated_at") values ('01J27JA3WM81T7D1TCNJ6Q06P5', '2024-09-04 22:24:43.686257+00', NULL, '01J6ZJ7PJR08D72KM5BAF67808', false, 'article', NULL, '20151204-php-neden-zamanin-ruhunu-yansitmiyor', 'published', NULL, NULL);
insert into "story" ("author_profile_id", "created_at", "deleted_at", "id", "is_featured", "kind", "properties", "slug", "status", "story_picture_uri", "updated_at") values ('01J27JA3WM81T7D1TCNJ6Q06P5', '2024-09-04 22:24:44.232534+00', NULL, '01J6ZJ7Q3VH05VHF4889WWE29V', false, 'article', NULL, '20160221-javascriptin-backend-gelecegi', 'published', NULL, NULL);
insert into "story" ("author_profile_id", "created_at", "deleted_at", "id", "is_featured", "kind", "properties", "slug", "status", "story_picture_uri", "updated_at") values ('01J27JA3WM81T7D1TCNJ6Q06P5', '2024-09-04 22:24:45.10021+00', NULL, '01J6ZJ7QYWGM91G0Q6HQRPPPES', false, 'article', NULL, '20111107-dotnetde-task-parallel-library-ile-producer-consumer-queue', 'published', NULL, NULL);
insert into "story" ("author_profile_id", "created_at", "deleted_at", "id", "is_featured", "kind", "properties", "slug", "status", "story_picture_uri", "updated_at") values ('01J27JA3WM81T7D1TCNJ6Q06P5', '2024-09-04 22:24:43.956782+00', NULL, '01J6ZJ7PV71HMMS0R23NCRMKGT', false, 'article', NULL, '20130428-standartlar-guidelinelar-ve-best-practiseler', 'published', NULL, NULL);
insert into "story" ("author_profile_id", "created_at", "deleted_at", "id", "is_featured", "kind", "properties", "slug", "status", "story_picture_uri", "updated_at") values ('01J27JA3WM81T7D1TCNJ6Q06P5', '2024-09-04 22:24:43.510896+00', NULL, '01J6ZJ7PD8J8YP66QE91Q2DH5Q', false, 'article', NULL, '20111023-dotnet-ile-gsmencoding', 'published', NULL, NULL);
insert into "story" ("author_profile_id", "created_at", "deleted_at", "id", "is_featured", "kind", "properties", "slug", "status", "story_picture_uri", "updated_at") values ('01J27JA3WM81T7D1TCNJ6Q06P5', '2024-09-04 22:24:44.924375+00', NULL, '01J6ZJ7QSGP0K3EDZM0CNCKN74', false, 'article', NULL, '20141219-eksik-parca-data-tanimlama', 'published', NULL, NULL);
insert into "story" ("author_profile_id", "created_at", "deleted_at", "id", "is_featured", "kind", "properties", "slug", "status", "story_picture_uri", "updated_at") values ('01J27JA3WM81T7D1TCNJ6Q06P5', '2024-09-04 22:24:42.313649+00', NULL, '01J6ZJ7N7VSV25K6ZD7ZVMQ0KV', false, 'article', NULL, '20221004-devrel-01', 'published', NULL, NULL);
insert into "story" ("author_profile_id", "created_at", "deleted_at", "id", "is_featured", "kind", "properties", "slug", "status", "story_picture_uri", "updated_at") values ('01J27JA3WM81T7D1TCNJ6Q06P5', '2024-09-04 22:24:45.373847+00', NULL, '01J6ZJ7R7GBWBFNK9ME20DJEVE', false, 'article', NULL, '20141122-tekerlegi-yeniden-kesfetmek', 'published', NULL, NULL);
insert into "story" ("author_profile_id", "created_at", "deleted_at", "id", "is_featured", "kind", "properties", "slug", "status", "story_picture_uri", "updated_at") values ('01J27JA3WM81T7D1TCNJ6Q06P5', '2024-09-04 22:24:44.654873+00', NULL, '01J6ZJ7QH103QMEQZK0QAQWJWY', false, 'article', NULL, '20141115-php-spl-standard-library', 'published', NULL, NULL);
insert into "story" ("author_profile_id", "created_at", "deleted_at", "id", "is_featured", "kind", "properties", "slug", "status", "story_picture_uri", "updated_at") values ('01J27JA3WM81T7D1TCNJ6Q06P5', '2024-09-04 22:24:42.593884+00', NULL, '01J6ZJ7NGMK53TAMW0C783FV0B', false, 'article', NULL, '20190205-gonullu-oldugum-islere-dair-planlama-ve-degerlendirme', 'published', NULL, NULL);
insert into "story" ("author_profile_id", "created_at", "deleted_at", "id", "is_featured", "kind", "properties", "slug", "status", "story_picture_uri", "updated_at") values ('01J27JA3WM81T7D1TCNJ6Q06P5', '2024-09-04 22:24:45.653402+00', NULL, '01J6ZJ7RG7RCMAFHWRX01BCDEM', false, 'article', NULL, '20130809-bilgisayarli-demokrasi', 'published', NULL, NULL);
insert into "story" ("author_profile_id", "created_at", "deleted_at", "id", "is_featured", "kind", "properties", "slug", "status", "story_picture_uri", "updated_at") values ('01J27JA3WM81T7D1TCNJ6Q06P5', '2024-09-04 22:24:41.627548+00', NULL, '01J6ZJ7MJDAG3GZHY49SE509RD', false, 'article', NULL, '20160610-microsoftun-yari-vizyonu-uzerinden-webcomponents-ve-spa', 'published', NULL, NULL);
insert into "story" ("author_profile_id", "created_at", "deleted_at", "id", "is_featured", "kind", "properties", "slug", "status", "story_picture_uri", "updated_at") values ('01J27JA3WM81T7D1TCNJ6Q06P5', '2024-09-04 22:24:43.819129+00', NULL, '01J6ZJ7PPXDJMG7H8M8K2C1BTF', false, 'article', NULL, '20160907-windows-komut-satiri-icin-kullanilabilecek-guclu-araclar', 'published', NULL, NULL);
insert into "story" ("author_profile_id", "created_at", "deleted_at", "id", "is_featured", "kind", "properties", "slug", "status", "story_picture_uri", "updated_at") values ('01J27JA3WM81T7D1TCNJ6Q06P5', '2024-09-04 22:24:46.687522+00', NULL, '01J6ZJ7SGH8HBPKXMWD3G0JGW9', false, 'article', NULL, '20161012-uzaktan-calisma-02', 'published', NULL, NULL);
insert into "story" ("author_profile_id", "created_at", "deleted_at", "id", "is_featured", "kind", "properties", "slug", "status", "story_picture_uri", "updated_at") values ('01J27JA3WM81T7D1TCNJ6Q06P5', '2024-09-04 22:24:46.225518+00', NULL, '01J6ZJ7S2487FJ2Q1MZB95ZYHE', false, 'article', NULL, '20230626-bilim-kurgu-distopya-tanisma', 'published', NULL, NULL);
insert into "story" ("author_profile_id", "created_at", "deleted_at", "id", "is_featured", "kind", "properties", "slug", "status", "story_picture_uri", "updated_at") values ('01J27JA3WM81T7D1TCNJ6Q06P5', '2024-09-04 22:24:44.385626+00', NULL, '01J6ZJ7Q8MYTANHEV0THKJ2XVN', false, 'article', NULL, '20221009-mentorluk-ve-kocluk', 'published', NULL, NULL);
insert into "story" ("author_profile_id", "created_at", "deleted_at", "id", "is_featured", "kind", "properties", "slug", "status", "story_picture_uri", "updated_at") values ('01J27JA3WM81T7D1TCNJ6Q06P5', '2024-09-04 22:24:47.112296+00', NULL, '01J6ZJ7SXTDJA9SZ9WGV7RD03H', false, 'article', NULL, '20181205-reactin-vizyonu-ve-hooks', 'published', NULL, NULL);
insert into "story" ("author_profile_id", "created_at", "deleted_at", "id", "is_featured", "kind", "properties", "slug", "status", "story_picture_uri", "updated_at") values ('01J27JA3WM81T7D1TCNJ6Q06P5', '2024-09-04 22:24:45.235234+00', NULL, '01J6ZJ7R34NJXYD7X854XCH5AG', false, 'article', NULL, '20181225-teslimat-bizim-isimiz-mi', 'published', NULL, NULL);
insert into "story" ("author_profile_id", "created_at", "deleted_at", "id", "is_featured", "kind", "properties", "slug", "status", "story_picture_uri", "updated_at") values ('01J27JA3WM81T7D1TCNJ6Q06P5', '2024-09-04 22:24:46.087131+00', NULL, '01J6ZJ7RXSVEDNK55KKZZ8VZ2N', false, 'article', NULL, '20191029-acik-kaynak-inisiyatifi', 'published', NULL, NULL);
insert into "story" ("author_profile_id", "created_at", "deleted_at", "id", "is_featured", "kind", "properties", "slug", "status", "story_picture_uri", "updated_at") values ('01J27JA3WM81T7D1TCNJ6Q06P5', '2024-09-04 22:24:41.760573+00', NULL, '01J6ZJ7MPJHY28R4M7KWF28SJE', false, 'article', NULL, '20161010-uzaktan-calisma-01', 'published', NULL, NULL);
insert into "story" ("author_profile_id", "created_at", "deleted_at", "id", "is_featured", "kind", "properties", "slug", "status", "story_picture_uri", "updated_at") values ('01J27JA3WM81T7D1TCNJ6Q06P5', '2024-09-04 22:24:45.79335+00', NULL, '01J6ZJ7RMK5Z7ETSNTHJYAVZ3Z', false, 'article', NULL, '20230217-acik-yazilim-agi', 'published', NULL, NULL);
insert into "story" ("author_profile_id", "created_at", "deleted_at", "id", "is_featured", "kind", "properties", "slug", "status", "story_picture_uri", "updated_at") values ('01J27JA3WM81T7D1TCNJ6Q06P5', '2024-09-04 22:24:42.174582+00', NULL, '01J6ZJ7N3G5HWRNDZ8YRJD8CDJ', false, 'article', NULL, '20160330-acik-kaynak-felsefesi', 'published', NULL, NULL);
insert into "story" ("author_profile_id", "created_at", "deleted_at", "id", "is_featured", "kind", "properties", "slug", "status", "story_picture_uri", "updated_at") values ('01J27JA3WM81T7D1TCNJ6Q06P5', '2024-09-04 22:24:43.014109+00', NULL, '01J6ZJ7NXRZC49DCQFFK8RV16R', false, 'article', NULL, '20090930-50-yil-icinde-tukenecekler', 'published', NULL, NULL);
insert into "story" ("author_profile_id", "created_at", "deleted_at", "id", "is_featured", "kind", "properties", "slug", "status", "story_picture_uri", "updated_at") values ('01J27JA3WM81T7D1TCNJ6Q06P5', '2024-09-04 22:24:41.469413+00', NULL, '01J6ZJ7KSEQ6ERWW48T3BZCQRY', false, 'article', NULL, '20160330-microsoftun-guclu-komut-satiri-perde-2', 'published', NULL, NULL);
insert into "story" ("author_profile_id", "created_at", "deleted_at", "id", "is_featured", "kind", "properties", "slug", "status", "story_picture_uri", "updated_at") values ('01J27JA3WM81T7D1TCNJ6Q06P5', '2024-09-04 22:24:46.4024+00', NULL, '01J6ZJ7S7MEQH7T95RG9NBBYFY', false, 'article', NULL, '20141122-icerik-sahipligine-dair-yasal-eksiklikler', 'published', NULL, NULL);
insert into "story" ("author_profile_id", "created_at", "deleted_at", "id", "is_featured", "kind", "properties", "slug", "status", "story_picture_uri", "updated_at") values ('01J27JA3WM81T7D1TCNJ6Q06P5', '2024-09-04 22:24:45.950561+00', NULL, '01J6ZJ7RSFJV6VBYPT467BKFQ0', false, 'article', NULL, '20221008-shu-ha-ri', 'published', NULL, NULL);
insert into "story" ("author_profile_id", "created_at", "deleted_at", "id", "is_featured", "kind", "properties", "slug", "status", "story_picture_uri", "updated_at") values ('01J27JA3WM81T7D1TCNJ6Q06P5', '2024-09-04 22:24:43.191998+00', NULL, '01J6ZJ7P393E92TMR12K3GY0FS', false, 'article', NULL, '20230826-iletisim-ilkelerim', 'published', NULL, NULL);
insert into "story" ("author_profile_id", "created_at", "deleted_at", "id", "is_featured", "kind", "properties", "slug", "status", "story_picture_uri", "updated_at") values ('01J27JA3WM81T7D1TCNJ6Q06P5', '2024-09-04 22:24:44.792592+00', NULL, '01J6ZJ7QNBGTHZPRTAN02SPJXQ', false, 'article', NULL, '20150103-microsofta-guclu-bir-komut-satiri-icin-acik-mektup', 'published', NULL, NULL);
insert into "story" ("author_profile_id", "created_at", "deleted_at", "id", "is_featured", "kind", "properties", "slug", "status", "story_picture_uri", "updated_at") values ('01J27JA3WM81T7D1TCNJ6Q06P5', '2024-07-29 07:36:36.768032+00', NULL, '01J3YPPRB2EWZK7N31DMC50Y8C', true, 'news', NULL, '20250518-hello-world', 'published', NULL, NULL);
insert into "story" ("author_profile_id", "created_at", "deleted_at", "id", "is_featured", "kind", "properties", "slug", "status", "story_picture_uri", "updated_at") values ('01J27JA3WM81T7D1TCNJ6Q06P5', '2024-12-03 14:18:21.353+00', NULL, '01JE63T76EBA8SPHPX9VJ9K1M0', false, 'article', NULL, '20241203-developer-platformu-eksikligi', 'published', 'https://objects.aya.is/story_20241203-developer-platformu-eksikligi.png', NULL);

insert into "story_tx" ("content", "locale_code", "story_id", "summary", "title") values ('Bu konuyla ilgili yazarken elimden geldiğince kimseyi incitmemeye çalışmaya özen
göstereceğim. Bazen aşırı uçta görüşlerimin olması nedeniyle, kendi ideallerimi
bu yazıya aktarsaydım bu yazıyı okunabilir bir şekilde sunmamın imkanı olmazdı.

Türkiyedeki açık kaynak topluluğu ile hiç dirsek temasınız oldu mu bilmiyorum
ben bir kaçının içinde, bir kaçına da yakınım. Elimden geldiğince faydalı olmaya
çalıştığım bir alan açık kaynak dünyası. Her ay -özellikle kendi kullandığım-
bir kaç projeye commit atarak kendimce katılımımı sağlamaya çalışıyorum.

Kimsenin benim katılımlarımdan haberi olmadığı gibi, bazen ben bile özellikle
takip ettiğim kişilerin çalışmalarından / katkılarından bi haber oluyorum.
Ülkede yapılan bir çok konferansın, konuşmanın ise “açık kaynak sevgisi”ne
rağmen “Microsoft ❤ Linux”den öte organik / çıkar odaklı ilişkiyi geçmediğini
gözlemliyorum.

Fakat bu konuşmaların hangi birinde “projelerinizde size yardımcı olacak yeni
ürün” dışında başka bir yazılımın mimarisinin konuşulduğunu veya ortak proje
geliştirmek/manifesto yayınlamak için ekip kurulduğunu gördünüz? Ben görmedim.
“Dünyada da böyle değil” argümanına saygım olsa da kendi evimin önü dururken
beni dünyadaki kirlilik ilgilendirmiyor açıkçası.

Ben yukarıdaki durumu görmediğim gibi aktif olarak başka projelere katılım
sağlayanların saygı gördüğünü, bu konferanslara konuşmacıyı geçelim konuk olarak
çağrıldığını da görmedim. Dolayısıyla o konferanslara katılımlarda her zaman
bize şeker olarak uzatılan “networking” kısmının daha çok “işsiz kalırsam bana
yararı dokunacak iş çevresi”nden ibaret olduğunu düşünmeye başladım.

Problem tanımını ortaya koyduğuma göre şimdi çözüm kısmına geçeyim,

Ortak bir platform oluşturarak;

- “katılım bekleyen” kendimizin veya topluluğun projeleri tanıttığımız,
- Bootstrap’e / PHP’e / Node.js’e v.s. nasıl commit göndeririz diye gereksinim
  analizi yaptığımız,
- Açık kaynak katılımıyla ilgili bilinmesi gerekenleri dokümante ettiğimiz,
- Kendi katılımlarımızı tanıttığımız, çözdüğümüz problemleri anlattığımız,

Bir ortam sağlamak. Bunun için en iyi ortam forum mu olur bilemedim. Ama eğer
yorumlar kısmında katılım gösterirseniz böyle bir çalışmaya en azından başlangıç
için girişmek istiyorum.
', 'tr          ', '01J6ZJ7QCTXKKS7DJ9NYMC9ZSH', '', 'Türkiyedeki açık kaynak topluluğu');
insert into "story_tx" ("content", "locale_code", "story_id", "summary", "title") values ('Geçtiğimiz hafta ODTÜ’nün KKTC kampüsünde gerçekleşen Cloud Computing and
Internet Security Seminer’ine katıldım. IBM, Tubitak ve Kaspersky’den gelen
konuşmacılar temsili olarak geldikleri kurumların Cloud Computing’e nasıl
baktıkları hakkında bende bazı izlenimler bıraktılar.

Cloud Computing Client-Server mimarisindeki “Server”ın monarşik rolünü tamamen
kaldırıyor. Ve devamında cloud avantajlarını kullanmak isteyen her
organizasyon/proje için hem yazılım hem de donanım olarak sıfırdan bir
yapılandırma gerektiriyor. Bu da bir çok IT firması için yeni fırsatlar anlamına
geldiğinden Cloud Computing’i IT sektörünün yeni para kazanma metodu olarak
değerlendirmeliyiz.

Her ne kadar servis üreten yazılım firmalarının güçlenmesi beni mesleki olarak
memnun edecek olsa da, bu senaryo benim gibi özgürlükçü birinin zihninde ister
istemez kaygılar oluşuyor. İşyerleri için Dedicated Hostinglerin, kişiler için
de External Disklerin azaldığı yerlerine verilerin servislere taşındığı bir
dünyadan bahsediyoruz. Başka bir değişle kişisel/kurumsal datanız üzerindeki
sahipliğinizi tamamiyle servisi aldığımız firmanın/organizasyonun sanal
storage’ı ile paylaşıyoruz. Bu firmalarınsa bizim haklarımıza ne kadar saygı
göstermesi gerektiğini ise ancak ve ancak geleneksel yasalar belirleyebiliyor.

“E yasa işte?” demek yeterli olmuyor. Facebook ve Google gibi firmalar şimdiden
kişisel bilgilerimizle bizi profilleyerek yasal yollardan olsa dahi bizim tercih
etmeyeceğimiz şekilde gizlilik haklarımızı çiğneyebiliyorlar. Şu an bir telefon
aldığınızı, arayanın mobil operatörünüzün olduğunu ve “x firması ile
bilgilerinizi paylaşmak isteyip istemediğini” varsayın, cevabınız %80 hayır
olurdu. Fakat social networkler ve çeşitli servisler bunu haberiniz olmadan
yapıyorlar.

Bir çözümüm olmasaydı bu yazıyı yazmayı da düşünmüyor olacaktım, işte benim
çözümüm. Cloud Computing, OpenSource Community’den ilham almalı ve
servis-kullanıcı ilişkisini lisans zorunluluğu ile yönetmeli. Kullanıcıların
ortak kullanacağı lisans tipleriyle şu an Apple’dan tutun Google’a kadar bizim
hakkımızda bilgi sahibi olan kurumların ellerindeki datamız üzerinde söz sahibi
olmamız sağlanabilir. Aklıma gelen örnek bir kaç kullanım şekli:

- Kişisel bilgileri servislerin kullanımına açmak,
- Kişisel bilgileri ancak servisde tanımlı kişilerin (kullanıcının kendisi,
  arkadaş listesi, v.s.) kullanımına açmak,
- Kişisel bilgilerden edinilecek gelirin kullanıcıyla paylaşılmak kaydıyla
  bilgileri servisin kullanımına açmak,
- Kişisel bilgileri yalnızca gelir amacı taşımayan organizasyonlarca kullanımına
  açmak.
- Kişisel bilgilerin servis tarafından kullanılmaması.
- Kişisel bilgilerin ve yedeklerinin servis üzerinden kullanıcının istediği an
  silinmesi.
- Kişisel bilgilerin servis üzerinden kullanıcının istediği an download
  edebilmesi veya başka bir servise taşınması.

“Taşınması” evet, madem Cloud Computing’den bahsediyoruz ben Gmail üzerindeki
tüm bilgimi Fastmail’e geçirmek için işlemleri başlatabilmeliyim. XMPP gibi
protokoller arası konuşma sağlayan ekstra protokollere Cloud sistemlerin
kullanıcılar ve abonelikler üzerinden konuşacağı varyasyonları olması gerekiyor
diye düşünüyorum. Aslında bu alanda prototip olarak bir çalışmam da mevcut, bir
ihtiyaç olduğu gerçek.
', 'tr          ', '01J6ZJ7NN6WHYB587G2TXRTS3Q', '', 'Cloud computing’de son kullanıcı açısından gizlilik');
insert into "story_tx" ("content", "locale_code", "story_id", "summary", "title") values ('Ben üretken olduğunu düşünen, ama asla ürettiklerimin yeterliliği konusunda
tatmin olmayan biriyim. Kendime bu konuda acımasız davranıyor ve yetersiz
görüyorum.

Ve daha başarılı olmak için işimden geriye kalan vakitlerde açık kaynak
projelere katılımda bulunmak, online toplulukları takip etmek ve yazılı/görsel
içerik üretme çabalarım oluyor. Bir süre sonra tekrar kendimi ölçmeye, geldiğim
noktayı değerlendirmeye çalışıyorum.

İşte bu değerlendirme esnasında kendi başarımı ölçümlemediğim tek metrik
“popülerlik”. Mesela bir blog yazısı özelinde düşünürsek okunma sayısı, youtube
videosu özelinde düşünürsek izlenme sayısı veya projeler için download sayısı.

Popülerliği değerlendirme dışında tutmamın nedenlerini paylaşırken, önce
genelgeçer durumu yorumlayıp, ardından kendi özelime geri döneceğim.

Bugün bir şeyler üreten bir mecra veya kişi maalesef ki tekil hit kavramından
kaçamıyor. Internet reklamcılığının dayattığı bu gerçeklik bizlere,
ulaşabildiğimiz, hitap ettiğimiz kitlenin geniş olmasını ve içeriğimizin
başlığını gören herkesin dikkatini bir sayfa kaydırma süresinde çekebilmemiz
gerektiğini telkin ediyor.

Üretebileceklerimizin çeşitlilik yelpazesini hedef kitlemizin mümkün olan en
büyük kesiminin ilgi alanını kapsayabilecek konulara indirgediğimizde,
birbirinin tekrarı yazıların, konferansların, videoların çoğalmasına tanıklık
ediyoruz. Neden mi çeşitlilik azalıyor, tekrarlar çoğalıyor? Herkesin damak
tadına uyacak bir lokantadan bahsediyoruz, çeşitleri elbette ki 5’i geçmeyen,
baharatsız yemekleri barındıracak.

Online içerik üretimi özelinde durum daha da kaygı verici. Ben BBS dönemlerinden
bu yana online içerik tüketicisi biri olarak, gün geçtikçe kaynaklarımın
azaldığını hissediyorum. Benim kişisel gözlemim, içerikler yavaş yavaş herhangi
bir derinlikli bilgi ağırlığı olmayan ürünlerin/markaların çevresinde
yoğunlaşmaya başlıyor. Örneğin “Separation of concerns” nedir bilmeyen bir
kitleye “Angular’a nasıl başlarız” içeriği servis ediliyor. Üstelik bu içerik
halihazırda en az 20 kere farklı bloglarda, youtube’da, kitaplarda yer alıyorken
üretiliyor.

Durumu şöyle örnekleyebilirim: Televizyonda ulusal kanalların öğle kuşağında
evlilik programları dışında izleyecek bir şey bulabilir misiniz? Bu programların
size hitap etmediğini düşünelim. Çoğunluğun izlemekten memnun olduğu “evlilik
programları” dışında bir seçeneğiniz yok. Burada bu programların varlığı değil,
tüm kanallarda başka bir şey olmaması eleştirimizin ana noktası.

Ben ise kendi halimde ürettiğim işlerde farklı bir yaklaşımı benimsemeye
çalışıyorum. Popülerlik ve hit benim gözümde kapitalist kavramlar. John Nash’ın
“oyun kuramı”nı hatırlayalım. Gerçekten birçoğumuzun gönüllü olarak yaptığı
üretimlerde topluluk adına maksimum faydayı bireysel başarılarla mı
kazanabiliriz?

Ben mutlaka bir hedefleme olacaksa daha rafine bir kitleyi hedeflemeyi, herkesin
kendi en iyi olduğu noktada farklılaşarak kendi deneyimlerini aktarmasının “daha
yararlı” olacağını düşünüyorum.

İşte bu yüzden nicel değil niteliksel ölçeklerle yapılacak değerlendirmeleri
dikkate almak gerektiğini düşünüyorum.

Eğer yazımı beğendiyseniz, beğen tuşuna basmayı, kanalıma abone olmayı
unutmayın. Oradaki sayılar artsın, 100.000 aboneye özel yazı yazmayı
düşünüyorum.
', 'tr          ', '01J6ZJ7MTZNZ207P1NYAA41A61', '', 'Üretmek ve Gönüllük');
insert into "story_tx" ("content", "locale_code", "story_id", "summary", "title") values ('Eminim WPF ile geliştirme yapan herkes native Win32 formları aradaki font
rendering farkının bilincindedir. TextBlock/Label objeleri bir resim ile
birlikte aynı container içerisinde render edilirken daha da flu bir hal
alıyorlar.

Microsoft bu durumun farkında, diğer yandan konuyla ilgili şikayetler de
artmaktayken eğer siz de bu durumdan müzdaripseniz .NET Framework’ün bir sonraki
sürümünü beklemek yerine
“[Give me back my ClearType](http://arbel.net/blog/archive/2007/02/02/give-me-back-my-cleartype.aspx)”
isimli yazıyı okuyarak sorunun üstesinden gelebilirsiniz.
', 'tr          ', '01J6ZJ7NCGNHWTP081Q3DD0BXN', '', 'WPF’de font rendering');
insert into "story_tx" ("content", "locale_code", "story_id", "summary", "title") values ('Lisansüstü eğitim “öğrenci olmak” pratiği olanlar için nispeten kolay, fakat
hızlı bitirmek istediğinizde neredeyse her dersin ciddi ödev/araştırmaları
altında boğulabileceğiniz bir süreçtir. Ben ufak bir hile yapıp (elbetteki
mevzuatta yer alan bir hakkımı kullanıp) ortalamamı yükseltmek için yazılım
mühendisliği ile ilgili ek bir seçiminde bulunmuştum.

O derste aldığım en yararlı bilgi 90larda amerikan yazılım endüstrisinde
production’a çıkan yazılım projesi oranının %60’larda olduğu ve bununla ilgili
sektörün nasıl refleksler geliştirdiği üzerine olmuştur.

Bugün bir istatistik olsa eminim ülkemizde hayal kırıklığı ile sonuçlanan
yazılım projesi oranı %60’dan yüksek çıkacaktır. Hem de önümüzde hali hazırda
başarılı olan modeller ve sistemler var iken.

Hepimizin platform, yazılım ve programlama dillerinden tutun, iş
arkadaşlarımıza, işverenlerimize hatta gelirimize kadar hayal kırıklıklarımız
mevcut.

Ben konuyu naçizane kendi üslubumca değerlendirmek isterim. Elbette ki burada
toplumsal tespitlerden çok profesyonel kalıplarda kalmayı tercih ediyorum.

Hayal kırıklığı yaşadığımız şeyi X’in yerine koyarak sormamız gereken birkaç
soru olduğuna inanıyorum.

- Bu hayal kırıklığı X’in karşılayacağını iddia ettiği ve karşılamadığı
  durumlardan mı dolayı mı oluştu?
- X’in başarısızlık riskine karşı organizasyonunuzun toleranslı olup olmadığı
  konusunda bir analiz yapmış mıydınız?
- X ona yüklediğiniz anlam için doğru tercih miydi? Farklı tercihler var mıydı?
  Bunların kıyaslaması sağlıklı yapıldı mı?
- Organizasyonunuz operasyonunuz için kusur barındırıyor muydu? X taşımaması
  gereken bir yükü taşıyor muydu?

Şimdi örnek senaryo üzerinden bir uyarlama yapalım ve sorulara bir kez de bu
şekilde bakalım. Diyelim ki Node.js ile yapılmış bir projeniz var ve Node.js’in
size yaşattığı deneyimler nedeniyle hayal kırıklığı yaşadınız.

- Node.js ile uygulamanızı yazdınız ve production’a çıktınız. Node.js’den
  kaynaklı bir memory leak ile karşılaştınız. Node.js’i tercih etmeden önce ön
  araştırma yapmış mıydınız? En azından “known issues”a göz geçirmiş miydiniz?

Şu anda uygulamanız çalışmaz halde bunun için alternatif bir planınız var mıydı?
Önlem olarak alabileceğiniz önceden planlanmış bir aksiyonunuz yok mu? Servisin
bir bölümünü geçici olarak kapatmak, kısa sürede workaround yazabilmek, v.s.?

- En başta tek seçeneğiniz Node.js miydi? Stres test ile şu anda yaşadığınız
  sorunun benzerini bu tercihi yaparken oluşturabilmiş miydiniz? Madem sıfır
  risk toleransınız vardı Node.js gibi nispeten yeni bir ekosistem yerine daha
  olgun bir platform seçemez miydiniz?
- Node.js yazabilen tek senior yazılımcınız şu anda problemi düzeltmekle mi
  ilgileniyor? Test yazmak, sistemi izlemek, continuous integration oluşturmak
  gibi DevOps operasyonlarını atlayarak mı bugünlere geldiniz?

Bu soruları çevrenizdeki insanlara, yeni gireceğiniz bir işe veya başka
kavramlara da uyarlayabileceğinizi düşünüyorum.
', 'tr          ', '01J6ZJ7MZ7DF5QEDZQFM7K8NV5', '', 'Hayal kırıklığına uğramamak için sorgulanması gerekenler');
insert into "story_tx" ("content", "locale_code", "story_id", "summary", "title") values (E'Ekşi Sözlük gibi mecralarda bir çok nostalji kuşağı estiren başlık bulunuyor.
“90larda çocuk olmak”, “magnum’un pahalı olduğu yıllar” ve benzerleri. Ben henüz
girilmemiş bir alana girip programlama özelinde kendi nostaljik deneyimimi,
90’larda programlama öğrenen çocuğun hikayesini paylaşmak istiyorum.

1993 yılında yurtdışından getirdiğim 486sx işlemcili masaüstü pc ile geçirdiğim
oyun dolu bir sene ertesinde Wolfenstein ve Volfied’in başı çektiği oyun
koleksiyonumdan “sıkılmıştım”.

Bilgisayarımı teslim alırken “ne yapacağız şimdi” diye sorduğum “teknik servis
yetkilisi abi”den öğrendiğim `cd` ve `dir` komutları ile MS-DOS 6.00’ın komut
satırından 420MB’lık diskimin içerisinde dolaşıyordum. Muhasebeci teyzemin bir
gün beni izlerken “bak EXE olanları çalıştırabiliyorsun” demesi ile
“aydınlanmıştım”.

`win` yazıp Windows’a girip paintbrush’da iki resim yaptıktan sonra yine sıkılıp
yeni bir şeyleri kovalamış, Creative Suite’den çıkan bookworm uygulamasını
açtığımda karşıma gelen “loading” ekranın beliren solucanı bilgisayar virüsü
zannederek bilgisayarımı fişten çekmiştim.

“Böyle olmuyor” dedim, çocukken ansiklopedilerine aşık bir çocuk olarak aldığım
bilgisayarı kullanmayı öğrenmeliydim. O eğlence kutusuna ait hevesim geçmeye
başlıyordu yavaş yavaş. Harçlık biriktirmeye, anneanne eli öpüp çeşitli
bilgisayar dergileri almaya başlamıştım. “Bilgisayar Pazarı”, “PC World” ve
“BYTE” dergilerindeki örnek kodları bilgisayarıma geçirip, üzerinde
değişiklikler yaparak önce Basic, sonra Pascal’da bir şeyler yapmak oldukça
ilgimi çekti.

Bir süre komut satırında Alarm, Piano ve Mastermind (Sayı tahmini) uygulamaları
porföyümü (!) oluşturmuştu. Pascal ile TSR bile yazmıştım.

Fakat başka programlama dilleri bazı konularda daha “yetenekli” olabiliyordu.
Benim için değişim kapıdaydı. “BYTE” dergisinde “C ile Windows gibi grafik
pencereler oluşturun” makalesini görmüştüm. 3 sayfa kod, “ben de pencere
çizersem Windows gibi bir şey yapabilirim” diye hayal etmeye başladım.

Bugün bir programlama dilinden diğerine nasıl geçiyorsunuz? İnternetten videolar
izleyerek mi? Ben bu işin benim çocukluğumda nasıl olduğunu anlatayım.

Disket kutumdan “harcanacak” 3–4 disketi formatladım ve bilgisayarcıya gittim.
“Bilgisayarcı abi” “Turbo C” bulabileceğini fakat o da başka bir arkadaşının
bilgisayarından alacağı için 3–4 gün sonra gelmem gerektiğini söyledi. O 3–4 gün
geçmek bilmedi, kodu tekrar tekrar makaleden okuyup kafamda değişiklikler
yapmaya başlamıştım. Gün geldiğinde, bilgisayarcı abinin ARJ’la disketlere
böldüğü Turbo C’yi bilgisayarımda `C:\\Programs\\Borland\\TurboC` benzeri bir
klasöre yerleştirdim. Ve C ile tanıştım.

Daha sonra yine dergiler sayesinde Windows üzerinde Visual Basic ve Visual
FoxPro ile programlamaya devam ettim. Ve hikaye devam etti, şu anda yazılım
geliştirmek artık “mesleğim”.

Özlemiyorum desem yalan söylemiş olurum. Bir hikayenin filmini izlemek yerine
kitabını okumak kadar farklı geliyor bugünle karşılaştırdığımda bana.
', 'tr          ', '01J6ZJ7P8W06V3WY2R70S53E8W', '', '90’larda programlama öğrenmek');
insert into "story_tx" ("content", "locale_code", "story_id", "summary", "title") values (E'90larda şu anda ceplerimiz dolduran dokunmatik ekranlı, 7/24 internete bağlı
akıllı telefonların herkeste olacağını hayal edebilir miydik? Star Trek gibi
vizyoner bilim kurgu serilerinin uzak gelecek tasfirlerini değerlendirmeye
katmazsak, bunu okumanın bir yolu olduğunu düşünmekteyim. O yolu paylaşıp,
“Yarın ne var?” sorusunun yanıtını aramaya başlamadan önce “Teknolojik
Arz/Talep” ve “Sentez Ürün” olarak isimlendirdiğim kavramları anlatmam gerekli.

#### Teknolojik Arz/Talep

Teknoloji hepimizin dilinde olan bir kavram. Hepimiz bu kelimenin kesin anlamını
bilmeden de olsa doğru bir şekilde kullanmaya “alıştık”. Gündelik konuşma
dilinde birçok kelime için bununla karşılaşabiliyoruz. Örneğin “apron”
kelimesinin tam olarak havalimanında nereyi kast ettiğini bilemeyebiliriz. Uçağa
uzanan körük müdür apron? Uçak pisti apron bölgesine dahil midir? Bu soruların
yanıtlarını bilmeden bu kelimeyi başarılı bir şekilde kullanabiliyoruz. Aynı
durum teknoloji kelimesinin anlamı için de geçerli.

**Peki nedir Teknoloji?**\\
TDK’ya göre “İnsanın maddi çevresini denetlemek ve değiştirmek amacıyla
geliştirdiği araç gereçlerle bunlara ilişkin bilgilerin tümü”. Ben
sadeleştirerek “insanların hayatlarını kolaylaştırmak için hayatına soktuğu
araçlar” diyorum. Günümüzde daha havalı şeylere “teknoloji” olarak baksak dahi
“tekerlek” de, “masa” da, “sandalye” de birer teknolojidir.

**Arz/Talep kısmı**\\
Arz ve Talep kavramlarını mutlaka bir yerlerden duymuşsunuzdur diye düşünüyorum.
Benim teorime göre “iş”leri sürekli kazanç elde etmek olan bilişim firmalarının
ne arz etmeye hazırlandığını okumanın geleceği okumakta önemli bir veri
oluşturuyor. Çünkü arz geldiğine göre,

- Artık bu teknoloji orta vadede kar edilebilir bir olgunluğa ulaştı,
- Kitleler tarafından benimsenmeye hazır bir teknolojiden bahsediyoruz,
- Bu teknolojinin arzında başarılı olan firma aslan payını kapacak, aynı
  sektörde olan firmalar öne geçmeden rekabette fark yaratmaları
  gerekiyor/bekleniyor,

**Talebin oluşumu**\\
Bahsettiğim firmalar elbette Microsoft, Apple, Google, IBM, Amazon ve Facebook
gibi devler. Bu firmalar konferanslarında etkileyici sunumlarla bizim
beklentilerimizi şekillendirecek, biz de mental olarak hazırlandığımız bu
kavramları tıpkı akıllı telefonlar, tabletler gibi hayatımıza kabullenmeye
başlayacağız.

Neden mi “bizim hazırlığımız” önemli? Microsoft’un 2000’lerin başında Tablet PC
denemesini hatırlayalım. Talep oluşmadan yapılan bir “kaza” arzıydı bana göre.
Tablet PC’den yaklaşık 10 yıl sonra daha estetik ve kullanılabilir bir ürünle
iPad başarıya ulaştı. Çünkü akıllı telefonlar, ipod, iphone derken bir zincirin
devam halkası oldu iPad. Kullanıcıları artık “dokunmatik ekranlı, kolay
taşınabilir, pc kadar güçlü olmayan ama iş için kullanılabilecek kadar akıllı
cihaz” teknolojisini hayatına almaya hazırlanmıştı.

#### Sentez Ürün

“Günümüzün akıllı telefonları nasıl meydana geldi?” sorusunun yanıtı bana göre
“birer birer”. Birbiri ile ilgisiz görünen aslında bu dönemin insanların temel
ihtiyaçlar bütününe seslenen ürünler önce bağımsız olarak karşımıza çıktılar.

- Akıllı telefon işletim sistemleri
- Ufak, taşınabilir dijital fotoğraf makineleri
- Yüzlerce şarkı depolamanızı sağlayan dijital müzik çalarlar
- Mobil internet
- Dokunmatik ve renkli LCD ekranlar

Bu ürünler tek bir gövdede hayat bulup ceplerimize sığdılar.

#### Yarın İçin Yatırımlar

Teknoloji devlerinin arz hazırlıklarını bir araya getirerek sentez ürünü tahmin
etmeye çalışmak “yarın ne var?” sorusu hakkında bir şeylerin kafamızda
canlanmasına yardımcı oluyor.

Ben veri olarak devler tarafından satın alınan firmalar ve onların aktif
çalışmaları ve çıkarılan ürün sayısı yüksek oranlarda artan alanları göz önünde
bulunduruyorum. Bu verileri ve trendleri sizinle paylaşmam gerekirse;

**Akıllı Asistanlar**

- IBM DeepBlue’dan sonra 2005’de geliştirmeye başladığı Watson’u yapay zeka
  çalışmalarının ortasına koyuyor.
- Cortana 2009’da Microsoft tarafından geliştirilmeye başlandı. Windows Mobile,
  Windows 10, Xbox One ve Android platformlarında varlığını sürdürüyor.
- Siri Apple tarafından 28 Nisan 2010’da satın alındı. Önce iPhone 4S, ardından
  Sierra güncellemesi ile macOS’un bir özelliği oldu.
- Google Now/Assistant 2011’de Google Voice Search’in geliştirilmiş sürümü
  olarak ortaya atıldı, 2012’de Android Jelly Bean’in demonstrasyonunda
  kullanıldığından beri Android üzerinde çalışmakta.
- Google’ın çatı şirketi (eski Google Inc.) Alphabet 2014’de DeepMind isimli
  yapay zeka firmasını satın aldı.
- Amazon ise Kasım 2014’de Amazon Alexa’yı duyurup Echo ile birlikte kullanıma
  sundu.
- Elon Musk Open AI projesi ile ev işlerine yardımcı olacak akıllı robotlar
  üretmeyi hedeflediğini açıkladı.
- IBM Watson başta olmak üzere yapay zekalara özel bir işlemci üretimi üzerinde
  çalıştığını duyurdu.
- IBM Watson’un sağlık alanındaki çalışmalarını güçlendirmek için Merge
  Healthcare’i satın aldı.
- Samsung 5 Ekim 2016 akıllı asistan Viv’i satın aldı.
- Facebook, Google, Microsoft, IBM ve Amazon “Partnership on AI” ismi altında
  yapay zeka alanında ortak çalışacaklarını duyurdular.

**Bot Platformları**

Microsoft bot framework’ü üzerinde 45 bin, Facebook ise Messenger platformu
üzerinden 34 bin geliştiriciye ulaştı.

Google ise API.ai firmasını satın alarak 60 bin geliştiriciyi transfer etmiş
oldu.

Yine Microsoft’un Flow ürünü ile bot geliştirilerine iş akışı (workflow)
platformu sağladığını not etmek gerekiyor.

(Kaynak: Webrazzi)

**Internet of Things**

IoT devrimi bir sürpriz olmaktan çok uzak. Artık hükümetlerin IoT’li bir dönem
için hazırlık yaptığı, ülkemizde dahi Arçelik’in HomeWhiz örneğinde olduğu gibi
beyaz eşya ve elektronik cihazlarımızın bir olup internete bağlı bilgisayar ve
akıllı telefonlarımızı azınlıkta bırakacağı günler için geri sayıyoruz. Bu
kesinlik içerisinde firmalar kendini konumlandırmaya devam ediyorlar.

Benim yakından takip etmeye çalıştığım Microsoft’un yeni JavaScript
yorumlayıcısı Chakra’nın IoT cihazları hedeflediğini açıklaması ve Amazon’un AWS
servislerinde IoT ile ilgili servisler sunması göze çarpan gelişmeler.

Satın almalara baktığımızda yalnızca 2016 itibariyle Microsoft Solair, Sony
Altair, Samsung Harman, Cisco Jasper, Intel Yogitech, Amazon ise NICE’i satın
almış bulunmakta.

**Giyilebilir Teknoloji**

Yazı halihazırda gayet uzun olduğu için arka arkaya vereceğim ürün isimlerinin
kendini açıklayacaklarını düşünüyorum.

- Google Glass
- (Microsoft) HoloLens
- Apple Watch
- (Facebook) Oculus Rift
- (Sony) Playstation VR

Bu listelediklerim yalnızca ürünleştirilmiş olanlar, üzerinde çalışılan projeler
ise çok daha fazla.

#### Yarın Ne Var?

Bunları göz önünde bulundurarak yarını okuyabilir miyiz? Okuyabiliyorsak, nasıl
çıktılar elde edebileceğiz? Ben bu sorulara ancak kendi vizyonum ve yorumum
doğrultusunda yanıtlar verebilirim.

Bence,

Yarının dünyası bizim bugün “bilgisayar” olarak tanımladığımız giriş / çıkış
cihazlarını yani ekranlardan, klavyeye, mouse’dan bilgisayar kasasına bir çok
kavramı yıkmaya hazırlanıyor.

Öncelikle kullanıcı ile bilgisayar arasındaki ilişki artık kullanıcının adım
adım komutlar vermesi yerine kişisel asistanların olaylar arasında bağlantılar
kurarak kullanıcının karşısına çıkmasına evrilmeye çalışıyor. Örneğin
kullanıcının annesinin doğum gününü hatırlatırken bir kadının o yaş aralığında
alabileceği hediyeler de kullanıcının karşısına getiriliyor.

Cihazlar etrafımızı sararken, biz bu cihazlarla ses ile iletişim kurmaya daha da
yaklaşıyoruz.

Hatta tüm teknolojileri bir arada düşündüğümüzde uçak bileti satan bir bot’la
iletişime geçen kişisel asistanınızın takviminize kaydettiğiniz bir seyahat için
periyodik fiyat alması, sokakta arttırılmış gerçeklik gözlüğü ile dolaşırken
kahve mağazalarının önündeki yeni nesil kioskların sizi tanıyıp size favori
kahvenize yakın içecekler önermesi bir 5 yıl öncesine oranla gerçekleşmeye daha
yakın duruyor.

Ben düşüncelerimi birer bilim kurgu filmi sahnesi gibi anlatıyor olsam da,
bunlar gündelik alışveriş ihtiyaçlarımıza dayanan ve aslında hayatımızın
içerisinde sürekli yaşadığımız senaryoların birer yeniden yorumlanışı.
Tahminlerimin isabetliliğini elbette zaman gösterecek.

Elimizdeki kesin bilgi ise bugün bir şeylerin başladığı. Büyük firmalar bu
şekilde veya daha farklı şekillerde benim yorumlarımdaki senaryoları
gerçekleştirebilecek teknolojilere yatırım yapmaya başladılar.

Peki bu gelecekte biz bir pozisyon alabiliyor muyuz? Bu da farklı bir yazının
konusu olsun.

Ufak bir destekleyici güncelleme olarak
[Prof. Dr. Oğuz Manas ile 1991’de gerçekleştirilen “2000’li yıllarda oturduğumuz yerden alışveriş yapacağız!” başlıklı söyleşi](https://www.egirisim.com/prof-dr-oguz-manas-2000li-yillarda-oturdugumuz-yerden-alisveris-yapacagiz/)yi
buraya eklemenin uygun olacağını düşündüm.
', 'tr          ', '01J6ZJ7SBVZ2NHN2BE00GT00ZF', '', 'Yarını okuyabilir miyiz?');
insert into "story_tx" ("content", "locale_code", "story_id", "summary", "title") values (E'Her ne kadar tooling konusunda yavaş kaldıklarını düşünsem de, .NET Core
cephesinde Microsoft’un oluşturduğu yol haritasında ilerlediğini hissedebilmeye
başladık.

Elimizdekiler:

- .NET Standard 2.0’ın frameworkler arasında uyumu arttırması,
- Visual Studio for Mac’in duyurulması,
- Yakında çıkacak Visual Studio 2017,
- Ve 1.0.1 sürümünde olan .NET Core

.NET Core beni fazlasıyla heyecanlandırıyor, Microsoft ismini ilk telafuz etmeye
başladığından bu yana sürece tanıklık etmeye çalışıyorum.

Platformu tam destekleyen bir IDE’ye sahip olup, bir iki proje geliştirdikten ve
mevcut NuGet paketlerimi .NET Core’a uyumlu hale getirdikten sonra üzerinde
ancak sağlıklı yorum yapabileceğimi düşünüyor, şimdilik fikirlerimi kendime
saklıyorum.

Fakat Microsoft’dan beklemediğim, beni şaşırtan bazı sakarlıklara denk geldiğim
için, bunlar da platform’un tasarımı ile ilgili olmadığından bu yazıyı
hazırladım. Sonuçta macOS üzerinde pek fazla .NET Core deneyimlemeye çalışan
kişi olmadığına inanıyorum.

**Problem 1:**\\
Visual Studio for Mac’i macOS Sierra’ma kurmayı denediğimde aşağıdaki mesajı
aldım:

Takıldığım bir başka nokta da Visual Studio for Mac’i Console Application
geliştirmek için kullanabilir miyim sorusuna yanıt alamamış olmak oldu. Ürünün
sitesinde yalnızca Mobile Development odaklı özellikleri vurgulanmış, en iyisi
başka bir zaman tekrar denemekti.

**Problem 2:**\\
Farklı bir zaman diliminde tekrar bir download denemesinde bulundum. Bu sefer
kurulumun çalışmasında bir sorun yoktu. Kurulum downloadlara başladı ve neden
olduğunu bilmediğim bir şekilde Xamarin bileşenlerini indirmeyi 2–3 kere
denedikten sonra vazgeçti.

Bu bir download hatası mıydı? Benim geçici bir bağlantı problemim miydi? (Ki
internette dolaşmaya devam edebiliyordum) Bilmiyorum. Fakat network kablom
çıkmış olsa dahi beni o an bir mesajla uyarıp “tekrar dene” seçeneği sunmasını
beklerdim. Hatalardan sonra tekrar kurulumu yeniden başlatabiliyorsunuz, ama bir
network hatası için çok büyük bir ceza.

**Problem 3:**\\
Kurulumu tamamlayabildikten sonra, hemen yeni bir Console Application
oluşturmayı ve hiçbir değişiklik yapmaksızın çalıştırmayı deniyorum ve
karşılaştığım durum:

**Problem 4:**\\
IDE’yi boşverip .NET Core sitesinden 1.0.1 sürümü kurulumuyla gelen pkg’i
indirip çalıştırma yoluna gittim.

Burada da bir bilgilendirme eksiği ile karşılaştım. .NET Core’un önceki
sürümlerine sahip isem nasıl bir işlem gerçekleştirmem gerektiğiyle ilgili bir
içeriğe erişemedim.

pkg’i kurup bash’de ilk projemi oluşturmak için adım attığımda aşağıdaki hatayı
aldım:

```sh
~$ dotnet -bash: dotnet: command not found
```

Hatayı kendim gidererek ilerleyebildim. Bu arada siz de başlıktaki hatayı aratıp
bu yazıya denk geldiyseniz, eksik path tanımını aşağıda olduğu gibi bir symbolic
link aracılığı ile atlatabilirsiniz:

```sh
ln -s /usr/local/share/dotnet/dotnet /usr/local/bin/
```

Alternatif olarak, `~/.bash_profile` dosyanıza aşağıdaki satırı
ekleyebilirsiniz:

```sh
export PATH=$PATH:/usr/local/share/dotnet
```

Sonuçta karşılaştığım bu hatalar beni ürünün henüz test senaryolarının pek
olgunlaşmadığı yönünde düşündürüyor.

Yine de bir kez daha .NET Core’dan oldukça umutlu olduğumu belirtmek
isterim. .NET sevdiğim bir platform, Visual Studio 2017’in çıkışı ve tooling’in
artmasıyla artık benim de üzerine yatırım yapabileceğim bir olgunluğa sahip
olacak. Sabırsızlıkla bekliyoruz :)
', 'tr          ', '01J6ZJ7RBVQ3GD7H1NHN0ZYBSY', '', 'Visual Studio for Mac deneyimlerim');
insert into "story_tx" ("content", "locale_code", "story_id", "summary", "title") values (E'Sabah twitter üzerinden “PHP’nin kötü kararlar aldığını düşünüyorum. 4–5 sene
önce hem modern hem de popüler bir dil olmayı başarıyordu PHP. bugün ise
başaramıyor.” yorumunu yapmıştım. Kapalı bir ifadeydi, haliyle nedenlerine yer
verilmemişti.

Biraz açmak istiyorum.

2010''ların başında PHP neler yapıyordu:

- MySQL paketi yetersiz olduğu için MySQLi’yi ortaya atıyordu. Transaction
  desteği geliyordu.
- PHP 5 ile object oriented yapısını yeniden ele almış PHP durak vermiyor, SPL
  ile dile yeni objeler katıyordu.
- PHP’nin en güçlü yanı kütüphanesiydi sürekli yeni fonksiyonlar ekliyorlar,
  bunlara özellikler getiriyorlardı.
- Extensionlar sayesinde genişlemeye tam destek veriyordu.
- PDO ile veritabanı katmanını soyutlaştırıyor. Mümkün olduğunca abstraction
  üzerinden yazılımcı problemlerini çözmeye çalışıyorlardı.
- Düzenli kodlar üretilebilsin diye en popüler framework’ü CodeIgniter olmasına
  rağmen topluluğa öncülük ediyor ve Namespaceleri ortaya atıyordu.

2015''lerde PHP neler yapıyor:

- Performans tarafında iyileştirmeler yapıyor.
- Dil sentaksını genişletmeye çalışıyor.
- Kod tekrarının daha az olması için anonymous method ve classları oyuna
  sokuyor.

Özetle PHP eskiden dikey büyüyerek platforma yeni olanaklar sunarken, şimdi daha
çok özelliklerini konsolide ederek muhafaza etme konumuna geçti.

Ben bu ivmedeki düşüklüye ve eski cesaretli "yapı sökümcü" ruhun kaybına
razı olamıyorum. Hem PHP''e internallarını öğrenecek kadar yatırim yaptım, hem de
gönülden sevdiğim bir üründen bahsediyorum. Paralelde diğer programlama
dillerinin de kullanıcısı olduğum için yeni proje başlangıçlarında elim diğer
platformlara kaymaya başladı yılların PHP kullanıcısı olmama rağmen.

Öncelikle mimari bir dönüşüm lazımdı. Shared hosting''lerin olduğu dünyada CGI
üzerinden PHP devrim yapmış olabilir. Ancak şu anda containerların,
fonksiyonların portatif hale geldiğı bir dünyada PHP yeni bir şeyler
söyleyemedi. Halen “PHP her request’de tüm kodu tekrardan interpreter’dan
geçirip işlem tamamlandığında bunu yok eden bir yapıya sahip”. Evet bu kulağa
çok deterministik gelse dahi ben programlama ihtiyaçları hızla artan biri olarak
long-serving processler/application serverlar üzerinden PHP ile koşan bir
runtime da tasarlamayı isterdim.

- Paket yönetimi ile ilgili sırtını tamamen composer’a dayıyor. Her requestte
  tekrardan dependency sınıfları teker teker aranıyor. Bunu kendi application
  server’ında çalışan birçok platform dille bütünleşik bir şekilde sundu. En
  garibanları Node.js bile NPM ile PHP''den daha ileride.
- Asenkron işlemler yapmak imkansız. (async ile threading farklı konulardır)
- Mikroservislerin, container bazlı uygulamaların üretildiği dünyada soket
  iletişimi nonblocking değil.
- Zamanında PDO, MySQLi için hem sınıf hem fonksiyon interface’i sunan PHP
  değilmiş gibi, ısrarla isimlendirme sorunu olan fonksiyonları statik
  kütüphanelere taşımıyor. Cesaret eksiği büyük sorun. strpos() -> Str::pos(),
  str\\_replace() -> Str::replace() gibi.
- Symfony, Laravel, Doctrine gibi en popüler frameworklerinin kullanmasına,
  topluluktan talep olmasına rağmen rağmen ısrarla dahili bir annotation desteği
  getirmiyor.
- Application context’i yok. Her requestte aynı kodlar neden bootstrap edip
  duruyor anlayamıyorum. Artık CGI’a mahkum değiliz.

Ben PHP 7’nin changeloglarına baktığımda performans iyileştirmelerine rağmen
halen CPU’nun ancak tek core’unu kullanabilen, async işlem yapamayan, modern
yazılım ihtiyaçlarını karşılamamaya başlayan, topluluğuna önderlik edemeyen bir
ürün görüyorum. 1999’dan beri yazdığım PHP için bunu 2 senedir her geçen gün
daha fazla hissediyor ve yeni projelere başka teknolojiler ile başlamaya
yaklaşmış biri olarak PHP’de bir şeylerin değişmeye başlaması gerektiğini
düşünüyorum.

Aksi takdirde shared web hostinglerde kaç tane wordpress kurulu olduğu ile gurur
duymakla yetineceğiz. Bunun yerine kariyerinin yeni başında olan insanların
PHP’yi tercih etmesi için geleceğe yönelik modern yaklaşımları benimsemek
gerekiyor.

RPCler, GraphQLler, socketlerden bahsettiğimiz dünyada, frontend tooling
yükselirken PHP bu cephelere girebiliyor mu? OOP tarafında .NET veya Java''ya
alternatif olmak istiyorsa alanının en hızlı, güvenilir ve iyi araçlara sahip
platformu olmak için atılımları var mı?

Karşı tezleri memnuniyetle dinlemeye hazırım.
', 'tr          ', '01J6ZJ7PJR08D72KM5BAF67808', '', 'PHP neden zamanın ruhunu yansıtmıyor?');
insert into "story_tx" ("content", "locale_code", "story_id", "summary", "title") values (E'Geçtiğimiz 8–9 sene içerisinde frontend development tarafında artan ihtiyaçlara
paralel olarak; backend’den aşina olduğumuz kod tasarım konseptlerinin
(patternler, separation concernler, frameworkler, v.s.) frontend tarafında da
konjonktürel olarak uygulanabilir hale gelmesiyle JavaScript neredeyse bir
yükseliş dönemine girmişti. Bu popülerlik başta ECMA olmak üzere herkesi
JavaScript üzerinden daha da üretken olma yoluna soktu.

İnanılmaz verimli bir kaç senenin ardından, tjholowaychuk’un Node.js’den Go’ya
geçişi ile birlikte JavaScript’in yakaladığı ivme duraksamaya girdi. En azından
“mükemmel” olduğu ilüzyonu ile dokunulmazlık zırhının eskisi kadar güçlü değildi
artık.

Bunda elbette ki “öncül”lerin Node.js ekosisteminde ürün geliştirip,
deneyimlerini yaşamaları ve daha sonra süreci değerlendirip “feedback” vermeye
başlamaları da etkili oldu.

Bugünlerde ise her gün twitter’da JavaScript hakkında değerlendirmeler okumaya
başladım. Twitter fiziksel limitleri nedeniyle bu konuların “derinlemesine”
konuşulabileceği bir mecra olmadığından bu değerlendirmeleri sağlıklı
bulmuyordum. Yine benzer Quora’nın derinliklerinde “neden JavaScript’le olmadı”
benzeri başlıklara sahip yazıları okuma fırsatı buldum.

Görünen o ki 2 yıl önce “MEAN Stack” diye pazarlanan JavaScript’in backend’deki
geleceğine dair genel kanı “olumsuz” olmaya doğru yol alıyor. Her genel kanı
gibi bunu da edinilen “sonuç”lar belirliyor.

Ben “genel kanı”ların “deneyimler karşılığında oluşan sağlıklı önyargı”ları
yansıttığını düşünsem de, bunları oluşturan duruma özel konuları incelemeden
hemen kabul etmek konusunda oldukça tutucuyumdur. Çünkü sonuçların (başarısız
veya başarılı girişimler) süreçten bağımsız olarak değerlendirildiğinde tek
başına bir veri niteliği taşımadığını düşünmekteyim. Tüm bir değerlendirme
yaptığımda ise “henüz” beni JavaScript’in backend geleceğinin pek parlak
olmadığı konusunda halen ikna edebilmiş bir görüş bulunmamakta.

Buradan hareketle JavaScript’in backend’deki durumuna ilişkin üretilen bazı
ortak argümanları buraya taşıyarak, kendi mesleki görüşümle yanıtlamaya
çalışacağım.

Bunu yaparken de her ne kadar JavaScript’in backend tarafında
“uygulanabilirliği”ni Node.js ile sınırlı görmesem de, öncül olduğu ve yazının
yazıldığı tarihteki en ciddiye alınır girişim olduğu için Node.js ekosistemi
özelinde değerlendirme yapacağım.

Yapmayacağım şey ise JavaScript ve Node.js’in güçlü yanlarından ve
avantajlarından bahsetmek olacak. Yazıyı okumadan önce konular hakkında zaten
bilginizin olduğu varsayımında bulunuyorum.

**İddia 1: Debug ve profiling zor**\\
PHP gibi bir application’a (apache, fpm, v.s.) modül olarak çalışan platformlara
göre Node.js’in bu konularda araç zenginliğinden dahi bahsedebiliriz. Hatta
belirli kod parçalarını ayırıp Node.js yorumlayıcı veya browser üzerinde
debug/profile edebilmek ufak da olsa bir avantaj sağladığını belirtmek
gerekiyor.

Java veya C# gibi dillerin IDE’lerindeki rahatlıktan bahsedemeyiz ama benim
kullandığım Visual Studio Code yaklaşık bir IDE desteği sağlamaya çalışıyor.

Eğer konu ürettiğimiz kodun değil de, Node.js’in runtime stabilitesinden
kaynaklı durumların profiling’i ise; bunun “stabilite” başlığında değerlendirmek
isterim.

**İddia 2: Stabilite problemleri mevcut**\\
Öncelikle _semantic versioning_ hayranı olduğumu belirtmek isterim. Kendimi
bildim bileli hiçbir şartta kendini henüz “1.0” olarak tanımlamamış bir ürünü de
production’da kullandığımı hatırlamıyorum.

Fakat
[Why did Koding switch from Node.js to Go?](https://www.quora.com/Why-did-Koding-switch-from-Node-js-to-Go)
gibi hikayeleri okudukça da birçok hikayenin aslında Node.js’in “henüz kendini
stabil ilan etmediği” günlere atıf ettiğini görüyorum.

_Yine de_ bu eleştirilerin bugünün sürümleri olan 4.x ve 5.x için
gerçekleştiklerini varsayıp yorumlamaya çalışalım. Burada stabilite problemleri
iddiasına cevaben “stabilite problemi kesinlikle yoktur” gibi bir yanıt
yazılamaz elbette. Ancak karşı argümanlar sunarak ilerlersek:

- Bugüne kadar Node.js’in mimarisinden dolayı düzelemeyecek bir durum duyduğumu
  hatırlamıyorum. (varsa lütfen düzeltin)

- Node.js yakın zamanda release cycle’ını LTS (long term support) sürümler
  üretecek şekilde güncellemiş görünüyor.

- Google Chrome başta olmak üzere browserlardan text editorlere kadar bir çok
  yazılımda kullanılan V8’in runtime hatalarının zaman içerisinde düzelmesi siz
  “kusursuz başka bir platform” arayışı başlatıp bitirene kadar geçen süreden
  daha kısa olacaktır.

- Başlıca şikayetler arasında duyduğum “memory leak” başınıza gelse dahi bunu
  tespit ettiğiniz anda workaroundlar sayesinde V8’e gelecek patch’den önce
  telafi edilebilir durumda.

**İddia 3: Bizim için hayal kırıklığı oldu**\\
Bu yorum oldukça subjektif. Fakat ben bununla çok karşılaştığım için
[Hayal kırıklığına uğramamak için sorulması gerekenler](./20160221-hayal-kirikligina-ugramamak-icin-sorgulanmasi-gerekenler)
yazımda konu olarak ele aldım.

Projenizin “bizim için hayal kırıklığı oldu” ile sonuçlanmasını istemiyorsanız;
önceden ödevinizi iyi yapmanız veya bir danışmanlık almanız şart. Çünkü “bugünün
mevcut sorunları” bugün başlanılacak bir projenin “hangi teknoloji ile
yapılacağını” değiştirebilecek kadar ciddi bir faktör.

Diyelim bu durum sizin için geçerli oldu ve Node.js’den kaynaklı bir riskten
ötürü JavaScript kullanmak beklentilerinizi karşılamıyor; ben bunu yine
JavaScript’in sunucu taraflı geleceğini belirleyen bir parametre olarak kabul
edemem.

Nedeni ise her platform tercihinin “proje başına” yapılmasının gerekliliğine
inananan biri olarak, production’a çıkacak ciddi yazılımlar geliştirirken
beklentilerin ve önceliklerin her projeye göre ciddi değişkenlikler gösterdiğine
inanıyor olmam. Evet Java, .NET, Python v.s. gibi sağlamlığını ispat edebilmiş
daha olgun alternatifler her zaman orada durmaktalar. Ama bu alternatiflerin de
yer yer bazı stresleri kaldıramadığını bizzat kendim deneyimledim.

Bu nedenle benim için bir platformun bir projede “bekleneni en iyi veren tercih”
olmaması, o platformun sonuna kadar üzerinin çizili olması anlamına gelmemekte.

**İddia 4: Geleceği muallak**\\
Danışmanlık verdiğim kurumlarla teknoloji seçimi yaparken her zaman göz önünde
bulundurduğum önemli bir parametre var. O da mevcut teknoloji üzerine farklı
kurumlar tarafından yatırılan “sermaye”.

Şu anda yazılım alanında Microsoft gibi devlerin planlarını JavaScript’e
yönlendirdiğini ve V8’e alternatif olacak bir JavaScript yorumlayıcı üzerinde
çalıştıklarını biliyoruz. Yine Microsoft giderek open sourcelaşma
yolunda .NET’den sonra en güvendiği platformu Node.js olarak belirlemiş olmalı
ki gerek Node.js’e verdiği IDE desteği, gerekse .NET’in çoklu platformda
çalışacak yeni ekosistemi .NET Core için yardımcı ürünlerin çoğunu Node.js
platformunun üyelerinden seçiyor.

Ek olarak .NET’e senelerini vermiş bir profesyonel olarak Microsoft’un her zaman
“yalnızca müşterisi” olduğumu düşündüm. Fakat Microsoft’dan bir yetkili ilk kez
benim yazılım geliştirici görüşümü almak için beni bir Skype görüşmesine davet
ettiğinde konu “Node.js hakkındaki görüşlerim”di.

PayPal’ın da yüzlerce “JavaScript mühendisi” yetiştirmekten övündüğünü, bu
ciddiyette bir firmanın JavaScript üzerine bolca yatırım yaptığını gözden
kaçırmamalıyız.

Örnekler elbette ki çoğaltılabilir. Fakat ana noktayı kaçırmamak gerekli.
Sermaye o kadar önemli bir parametre ki, o ekosistemden beslenen büyük firma
sayısı arttığı sürece Node.js’in sunucu taraflı performansı ne kadar kötüleşirse
kötüleşsin bu durum birilerini farklı bir çözümle ürün devamlılığını sağlamak
zorunda bırakacaktır.

Node.js’i önemli kılan ortaya çıkışındaki JavaScript dilinin gücünün backend’de
kullanılabilecek olduğu teziydi. Node.js’in kendisi başarısızlığa doğru gitse
dahi tezinin geçerliliğini yitirmeyeceğini düşünüyorum.

Zaten hali hazırda browserlardaki tahtı 20 yıldır sarsılmamış bir JavaScript’den
bahsediyoruz. V8, SpiderMonkey, Carakan v.s. derken yakın gelecekte de
JavaScript yorumlayıcılarının her geçen gün gelişeceği kolaylıkla tahmin
edilebilir. Node.js’in da avantajı bunlar için bir çatı görevi sağlaması.
Örneğin geçtiğimiz günlerde Microsoft’un Chakra isimli JavaScript yorumlayıcısı
Node.js’e yapılan bir Pull Request ile artık V8 yerine alternatif olarak
kullanılabilir hale geldi. Şu anda Node.js iki farklı yorumlayıcı ile
çalışabiliyor.

**İddia 5: Yazdığınız kod gelecek sene çalışmayacak**\\
JavaScript’in standartları ECMA tarafından belirlenmekte. Çok büyük bir sorun
oluşturmadığı sürece geriye doğru uyumluluk sağladıklarını biliyoruz. Burada
browser API’lerinden bahsetmiyorsak zaten JavaScript’in uğradığı değişiklikler
çoğunlukla dile yeni yetenekler kazandırmaya yönelik.

Node.js ise _semantic versioning_ sayesinde API’sinde olan bir değişiklikte
sürüm değiştireceğini garantilemiş durumda. **NVM** gibi araçlar sayesinde aynı
makine üzerinde farklı Node.js sürümleri de kullanabiliyorsunuz. Ben burada bir
problem göremiyorum?

**İddia 6: Büyük mimariler için çok kötü**\\
Aynı şeyi PHP ve C/C++ için de söyleyebiliriz. Fakat aynı dillerin kendi
platformlarında en büyük codebase’lere ait projeler tarafından da kullanıldığını
biliyoruz. Peki bu nasıl oluyor?

Kod Organizasyonu! Node.js’in mevcut ekosistemindeki frameworkler monolitik
yapıların antitezi gibi duruyor. Bunun için kod organizasyonu ve dizin yapısını
oluşturmak tamamiyle size kalıyor. Eğer bu konuda iyi bir tasarım
çıkartamazsanız bu dertten muzdarip oluyorsunuz.

Belki de yapılan en büyük hata codebase’in sürümlenmemesi ve tasarımı konusunda
yeteri profesyonel yaklaşımların gösterilmemesi olabilir.

**İddia 7: Bakımı Zor**\\
Bakımın hangi açıdan zor olduğuna göre değişebilecek bir argümandan
bahsediyoruz.

Eğer yeni özelliklerin eklenebilmesinden, kodun sürdürülebilirliğinden
bahsediyorsak bu yine bir organizasyon sorunu. Coding Guidelineları takip etmez,
kod organizasyonunu eksik yapar ve dökümantasyona önem verilmezse her platformda
karşımıza çıkabilecek bir problemden bahsediyoruz.

Şayet sistem tarafında bir bakımdan söz ediyorsak Java’dan ve .NET’den sonra
çalıştığım en fazla araca sahip ekosisteme sahip olduğunu, ve araç sayısına
bakarsak en yüksek ivmeye sahip platform olduğundan rahatlıkla bahsedebiliriz.

**İddia 8: Hype**\\
“Node.js çevresinde oluşan hype’dan dolayı oldukça yaygınlaştı ve dikkatleri
çekti.” diyorsanız bir sonraki cümlenizde Node.js’in yerine ne koyduğunuzu
açıklarken kendinizi izleyin. Eğer bu Python, Java v.s. değil de Go ise; Go’nun
da bir hype ile bu noktaya gelmediğini iddia edemeyeceğinizi düşünüyorum. Burada
“Go” yerine pazarlaması heyecan yaratmış, sizin de kod yazdıkça öğrenmenin
mutluluğuyla çok sevdiğiniz fakat 5–6 proje sonra sıkılacağınız bir dili
koyabilirsiniz.

JavaScript’e dönersek, hype’dan beslenmiyor demek ne kadar isabetli olur
bilemiyorum ama client tarafındaki uzun yıllar sürdürdüğü popülerliğini göz
önünde bulundurduğumuzda bunun gelip geçici olduğunu iddia etmek oldukça zor.

Fatih Kadir Akın’ın bu yazıya cevaben
[Medium’da yer alan yazısı](https://medium.com/@fkadev/node-js-in-back-end-gelece%C4%9Fi-f246db981721)
üzerine bu noktadan itibaren yeni iddiaları ekliyorum.

**İddia 9: Event Loop üzerine kurulu ve bu kötü bir şey**\\
10 yıl önce Java’da veya .NET’de “ben web sayfası değil, backend development
yapıyorum” diyebilmek için platformların paralel/çoklu işlem kabiliyetlerini
kullanabilmek önemli bir eşikti benim için. GSM operatörleri ile sürekli soket
üzerinden iletişimde olan, SMS/MMS toplu mesaj gönderimi yapan ve son
kullanıcılara API sunan bir SaaS’da lead architect görevindeydim. Sabahlara
kadar Thread optimizasyonu yaptığımı, en ufak yanlış hesapta sistem
performansının anlık 5 kişiye zar-zor yanıt verebilecek hale geldiğini
hatırlıyorum o günlerden.

.NET, Java v.b. platformlar bu gücü size kullandırabildiği için “güçlü bir
platform” olarak tanımlanıyor. Fakat kağıt üzerinde basit olmasına rağmen
antremanlı ve deneyimli değilseniz forklar, joinler ve senkronizasyon
problemleri bilin ki sizi üzmeye geliyorlar.

Paralel/çoklu işlemleri gerçekleştirebilmek için en çok kullanılan mimarileri
incelediğimizde:

- Multi-threading / Multi-process yapılar

- Event-driven yapılar

- Hibrit yapılar

ile karşılaşacağız.

Bunlardan ilki gerçekten zahmet ve kontrol gerektiren bir yapı. Donanıma göre
thread/process sayısını ayarlamak, elimizdeki işleme göre dinamik olarak bir
thread havuzu mu yoksa sabit thread sayısı mı işimizi görecekse onun
implementasyonu ile ilgilenmek, kullanılan değişken tanımlarının senkronizasyonu
gibi faktörler yakanızı bırakmamakta.

İkincisi ise arayüzler, oyunlar ve komut satırlarının vazgeçilmezi ana/olay
döngüsü etrafında sağlanıyor. İnanın kodda bu ana döngü gerçekten “çok çirkin”
duruyor fakat performans olarak Thread’den pek de geride kalmayan, Thread’in de
yukarıda saydığım dezavantajlarını barındırmayan bir yapı.

Üçüncüsü madde yani hibrit sistemler ise bu yazının konusu değil.

JavaScript burada belli ki bir tercihte bulunmuş. Nasıl .NET’e, Java’ya,
Python’a pointerları yönetme gücünü bize bırakmadığı için kızmıyoruz,
JavaScript’in de “JavaScript yazılımcısının bu kadar low-level optimizasyonlarla
işi olmamalı” kararına hem katılıyor hem de saygı gösteriyorum. Node.js gibi
single-thread üzerinden olmasa da event-driven mimariyi nginx’in, lighttpd’nin,
tornado’nun, netty’nin de benimsediğini eklemek isterim. Dolayısıyla “event loop
öcü değil” diyebiliriz.

Yarın öbür gün Thread’ler kolaylıkla JavaScript’e getirilebilir ve platformdaki
yerini alabilir. Node.js de .NET runtime’ı gibi hibrit bir yapıya doğru geçiş
yapar, teknik olarak mümkün bu. Fakat Thread gibi mekanizmaların yanlış
kullanımı (ki çokça rastlanır) performanstan çok performans kaybı olarak
yazılımcıya geri döner. Bunları bazen tam olarak fark edemezsiniz bile.

**İddia 10: Backend için nonblocking IO kullanışlı değil**\\
[Değil mi?](http://stackoverflow.com/questions/10570246/what-is-non-blocking-or-asynchronous-i-o-in-node-js)
Ben aksine, özellikle yukarıda da bahsetmiş olduğum
mesajlaşma — kuyruklama — işleme — iletme işlemleri için tek bir CPU core’unun
ve RAM veriminin maksimize edilerek kullanılmasının oldukça önemli olduğunu
düşünmekteyim. Hardware’da pipelining’den tutalım birden fazla core içeren
işlemciler kullanmamıza kadar elimizdeki her teknoloji işlemlerin nonblocking
gerçekleşmesini hedef almışken biz işin software kısmında neden nonblocking
IO’yu es geçiyoruz ki?

.NET 1.0’dan beri içinde Threading bulunmasına rağmen, 4.0 sürümü ile gelen
“Task Parallel Library”nin pazarlamasını paralelizm üzerinden yapmıştı.

Dolayısıyla nonblocking IO’nun backend için yüksek işlem ölçeğine gidildikçe
öneminin arttığını düşünmekteyim.

**İddia 11: Callback Hell’den Promise’lere oradan Async/Await’e geldik**\\
Evet, iddia doğru olmasına rağmen bunun eleştirilmesini sıkıntılı buluyorum. En
azından bir dilin gelişip olgunlaşmasına karşı çıkacak durumda olduğumuzu
düşünmüyorum.

.NET’deki async/await’in JavaScript’e getirilmesini “tekerleği yeniden
keşfetmek” olarak tanımlamaktan çok, başarılı bir modelin aynen uygulanması
olduğunu düşünmekteyim.

.NET’de 1.1’den itibaren yer alan AsyncCallback/IAsyncResult tipleri
JavaScript’deki Promise’e benzer görevde bir nesne tipiydi. 2012’de .NET 4.5 ile
birlikte vitrindeki yerini alan async/await’in de JavaScript tarafından
benimsenmesiyle bana göre _Callback Hell_ olarak tanımladığımız sorundan
kurtulmuş olduk.

Ayrıca eklenmesi gereken bir not da .NET’in async/await’in kazanımından sonra
kendi kütüphanesindeki metodları mümkün olduğunca async olarak yazmaya
başladığıdır.

JavaScript’in kendi içerisinde yaşadığı bir problemi başarılı bir uygulamayı
benimseyerek aşıyor olması bana kötü değil, bilakis iyi bir strateji olarak
geliyor.

**İddia 12: Node.js sezgisel değil**\\
Genellikle bu eleştiriler Node.js API üzerinden şekil aldığı için en hak
verdiğim iddia bu.

Yine de savunma makamı görevimi bozmaksızın API’ın her zaman değişebileceğini,
maskelenebileceğini düşünmekteyim.

Özellikle Fatih Kadir Akın’ın yazısında callback üzerinden yaptığı eleştiriler
Node.js tarafından kafa yorulduğunda aşılamayacak bir durum değil. Başka bir
deyişle henüz bir PHP’den bahsetmiyoruz.

Örneğin ileride Node.js API’ının callback parametresi almayan kullanımları
Promise döndürür, biz de bunlara da bir `await` ekleriz böylece herkes memnun
olur.

```javascript
let file = await fs.readFile(“./a.txt”); console.log(file) //=> beklediğim şey.
```

Yine `SetTimeout` yerine bir `await delay(ms)` ile callback’den kurtulabilmemiz
mümkün.

```javascript
console.log(1);
await sleep(2000);
console.log(2);
console.log(3);
await sleep(1000);
console.log(4); // buradan devam ediyoruz :) promise.resolve’a da gerek yok :)
```

Böylece ben de async/await’e kullanım örnekleri vermiş oldum.

**İddia 13: Konvensiyon kötü, kod düzeni mevcut değil**\\
“Büyük mimariler için çok kötü” iddiasında da kod organizasyonu demiştim. Bu da
benzer bir konu. Yazılımı tasarlayan kişiye bağlı olarak sanat eseri de, çöp de
çıkartabileceğiniz bir malzemeden söz ediyoruz.

.NET ve Java’nın “framework” odaklı platformlar olduğunu unutmayalım. Bu
platformlar aynı zamanda uyulması gereken dizin yapısını da beraberinde
getiriyor.

Fakat C/C++/PHP’de nasıl sabit bir yapıdan ve konvensiyondan bahsedemezsek
JavaScript için de bu sabit yapıdan bahsedemiyoruz. Fakat bu alana dil girmedi
diye dil veya Node.js eleştirilmemeli. Bu durum aksine bize bir özgürlük
sağlamakta. Kısacası isteyenin istediği konvensiyonu takip edebilmesi
konvensiyon olmaması anlamına gelmiyor.

Bu özgürlüğün bir avantajı da 50–100 satırlık bir .NET projesini önerdiği
konvensiyon nedeniyle tek bir parça haline getiremezken, aynı JavaScript kodunu
Amazon Web Services’da lambda’ya atıp mikroservise dönüştürebiliyorsunuz.

**İddia 14: Çok fazla transpiler karmaşası var**\\
Evet, kesinlikle! Fakat bunu JavaScript’in ancak bugünlerine özel bir durum
olarak okumak gerekli. Oldukça büyük, ihtiyaca yönelik değişiklikler bir anda
dile eklendi ve V8’i geliştiren ekip buna adapte olmakta sıkıntılar yaşadı. Ben
önümüzdeki sene bu konuyu tekrar konuşuyor olacağımızı düşünmüyorum.

**İddia 15: Statik type’lı diller daha olgun ve güçlü**\\
Gereksinimler ve bunların tespit edilmesi işimizde büyük önem sahibi. Diğer
yandan yazılım ekibinin bilgi birikimi, yeteneklerini hangi dilde daha net
ortaya koyabildiği gibi pratik parametreler de tercihlerimizi etkileyen
faktörler.

Ben kodlamaya 90larda başlamış biri olarak;

- Pointerların, memory yönetiminin asla otomatik yönetilmemesi gerektiğini,

- Interpreterlardan programlama dili değil ancak script olabileceğini,

- Managed platformların asla başarıya ulaşamayacağını,

iddialarını bolca duydum. 2000’li yıllarda yukarıdaki önyargıların bir çoğunun
çürüdüğünü söylemek yanlış olmaz.

Ben 2010’lu yıllarda da strictly-typed yaklaşımların “pointer”lar gibi geride
bırakılacağını düşünmekteyim. Tabii ki type-hinting yaşayacaktır ama
runtimeların konsantrasyonu dynamic type yorumlamasına daha da yanaşacaktır.

Big data işlemekten bahsedip JSON okumak için 30 satır kod, 20 explicit type
casting’de bulunulduğumuzda sezgisellik de konvensiyon da geri planda kalmaya
başlıyor.

İşte bu karmaşada benim “kontrolde vites arttırma” olarak ifade ettiğim bir
durum gerçekleşiyor.

Eğer yazacağımız uygulama için Node.js de, Java da, C++ da uygunsa,

Benim ilk tercihim mümkün oldukça -zaman ve efor tasarrufu kaygısıyla- Node.js
kullanmaktan yana olur.

Daha çok mu kontrol gereksinimi duyuyorum, Thread yönetimi mi yapacağım.
Node.js’in üzerini çizip Java’yı gözüme kestiririm.

Yine mi tatmin olmadım? Pointer’a kadar inmem mi gerekiyor? Streamler üzerinde
mi çalışacağım? C++ veya C ile yoluma devam edebilirim.

Özetle daha olgun platformlar var evet, ama Node.js’in de avantaj olarak bize
sunduğu gerçekler var.

Sunduğum karşı fikirleri çürütecek itirazlar geldikçe veya yeni argümanlar
eklendikçe bu yazıyı geliştirmeyi planlamaktayım. Onun için yorumları
memnuniyetle yanıtlamak isterim :)
', 'tr          ', '01J6ZJ7Q3VH05VHF4889WWE29V', '', 'JavaScript’in backend geleceği');
insert into "story_tx" ("content", "locale_code", "story_id", "summary", "title") values ('Diskden bir klasör içerisindeki tüm dosyaları okumak, elimizde bulunan birden
fazla URL adresinden feedleri veya web sayfalarını download etmek ve elimizdeki
verilerin güncel kopyalarının olup olmadığını kontrol etmek. Bütün bu örnekler
için genellikle bir döngü yapısı kullanırız. Fakat kodu yazma amacımız belirli
bir sırayı takip etmesi değil de, tüm sistem kaynaklarını kullanarak işlemi
gerçekleştirmesi olduğunda birden fazla thread ile işlem yapmak en uygun
çözümdür.

Threading’le ilgili biraz bilginiz varsa “Producer/Consumer Queue” kavramını
duymuşsunuzdur. Bu kavramı en basit haliyle “dış bir kaynaktan beslenen iş
kuyruğu, ve bu iş kuyruğundan objeleri çeken birden fazla thread” olarak
tanımlayabiliriz.

Klasik thread yaklaşımında kullanacağımız kadar thread oluşturur ve bu
threadlerin sürekli veya iş kuyruğu tamamiyle boşalana kadar kuyruktaki objeleri
aldığı bir yapı oluştururuz. Yönetimi low-level bilgi gereksinimi istediğinden
zor ve mantık hatalarına açık bir yöntemdir.

.NET CLR bu nedenle .NET 2.0 ile birlikte ThreadPool yapısı, ardından .NET 4.0
ThreadPool üzerine oturttuğu daha high-level bir yapı olan TPL yani Task
Parallel Library ile oluşabilecek sorunlar konusunda yazılımcılara çözümler
sunmakta. ThreadPool ve TPL’in çalışma mekaniği yapılacak her kalem iş’e bir
Task (Görev) olarak bakmak, ve bu işleri yazılımcının yalnızca tanımlayıp
kodların işletim kısmıyla ilgilenmemesini, ancak bu görevleri iptal etmek veya
görevler tamamlandığında haberdar olmak gibi genel kısımlarda yazılımcının kod
yazmasını sağlar.

Fakat .NET 4.0 ile gelip, işlemcinin sanal/fiziksel tüm çekirdeklerini sonuna
kadar kullanmamızı sağlayan TPL’in bir kaç eksikliğinin olduğunu fark ettim.
Çalışacak maximum thread sayısına biz karar veremiyoruz, dolayısıyla vermiş
olduğumuz her kalem iş .NET tarafından içeride bir yerlerde kuyruğa dahil
edildiğinden kod ile bu kuyruğu izleyemiyoruz.

Bu nedenle kendi yazdığım ve aslında .NET’de oluşturduğum bir kütüphane olan
Tasslehoff’un parçası olan `ChannelTaskQueue`’nun kodlarını GitHub üzerinden
paylaşıyorum:

%[https://github.com/eser/tasslehoff-library/tree/master/Threading]
', 'tr          ', '01J6ZJ7QYWGM91G0Q6HQRPPPES', '', '.NET’de Task Parallel Library ile Producer/Consumer Queue');
insert into "story_tx" ("content", "locale_code", "story_id", "summary", "title") values ('Yazılım firmaları her ne kadar hizmet/çözüm sağlıyor olsalar dahi, sektör olarak
somut çıktımız yalnızca kod üretiminden ibaret. Bu üretim süreci boyunca code
review’ından, configuration management’ına kadar bir dizi araç kullanıyoruz.

Ben kişisel bir tercih olarak Crucible gibi araçlardan aldığım yardımın yanısıra
iki haftada bir üretilen kodu hardcopy olarak inceliyorum. Ekip arkadaşlarımın
ve açık kaynak olarak elime ulaşan bu kodları incelemek beni bir anlamda
eğitiyor.

Bu sürekli eğitim sürecinin bana öğrettikleri arasında önemli bir payı şu
çıkarım oluşturuyor:

- Olası bugların, regressionların, uyum sorunlarının büyük bir çoğunluğunun
  standartlar ve guidelinelara uyulmadığı noktalarda,
- Performans kayıplarının ve anlaşılmaz kodların ise best practise’lerin
  uygulanmadığı senaryolarda,

meydana gelmekte.

Bu nedenle genel kanının aksine yazılım hatalarının dikkatsizlikten değil,
araştırma ve eğitim eksikliğinden kaynaklandığını düşünmekteyim.

Eğitim aldığım iki farklı üniversitenin yazılım bölümlerinin çatısı altında
standartların anlatıldığına ve öneminin vurgulandığına denk gelmemiş olmam basit
bir rastlantı olmamalı.

Mevcut durum hakkında elbette derin study case’ler yapılarak incelemek daha
doğru olur fakat ben burada kendi perspektifimden yorumlarımı paylaşarak da bir
noktaya dikkat çekebileceğimi düşündüm.

Öncelikle yukarıda kullandığım kavramlar arasındaki farkları açıklamak
gerekirse:

**Standart**, özellikle platformlar ve kültürler arası iletişimler — aktarımlar
için kullanılır. Örneğin Unicode standartı sayesinde bugün Facebook gibi bir
platformda kril alfabesiyle yazılmış bir yazıyı ekranımızda okuyabiliyoruz.
Unicode’dan önce browser’ın dil ayarlarından bir kaç ayar değişikliği yapmamız
gerekirdi, bunun sonucunda da artık alfabemizde kullandığımız “ğüşöçı” gibi
karakterleri görememeye başlayacaktık. Kısacası siz unicode gibi bir standartı
yazılım içerisinde kullanmadığınızda yazdığınız iletişim formu japonca destek
sayfası olan bir firma tarafından kullanılamayacak.

**Guideline** ise bu kadar geniş çaplı bir mesele değil, bulunduğunuz ekipçe
proje başlarken gerek requirementlara gerekse keyfinize göre karar
verebileceğiniz bir dizi kurallar bütünü. Kod yazım kurallarından tutun,
veritabanında her kaydın ne zaman değiştirildiğine dair bir field kullanıp
kullanmayacağınıza yönelik kararlardan bahsediyorum.

**Best practise**’i ise bir işi yapmanın en doğru yolu olarak özetleyebilirim.
Tüm kullanıcıları veritabanından liste olarak çektiniz, bu kullanıcıların her
birinin aynı zamnada kaç yorumda bulunduğunu yazmak istiyorsunuz fakat
veritabanı veya veri depo yapınız basit bir JOIN ile işinizi görmeyecek şekilde
tasarlanmış. 20 kullanıcının listelendiği döngüde 20 kere ayrı sorgu yapmaktansa
olayın en başında 20 kullanıcının kimlik bilgilerini alıp, bu bilgilerle tek bir
seferde yorumları çekip kullanıcılara pay etmek çoğu zaman best practise’dir.

Benim naçizane görüşüm çalışmalarımızda bu kavramlarla “aşina” olmak yaptığımız
işin kalitesini yükseltecektir.
', 'tr          ', '01J6ZJ7PV71HMMS0R23NCRMKGT', '', 'Standartlar, guidelinelar ve best practiseler');
insert into "story_tx" ("content", "locale_code", "story_id", "summary", "title") values ('Cep telefonunlarımıza gelen kısa mesajlarda genel kanının aksine geleneksel
8bit’lik ASCII tablosu değil de, GSM operatörü tarafından özelleştirilebilen
7bit’lik GSM alfabesi kullanılmaktadır.

Bunun zorluğu ile en çok servis sağlayıcı tarafından filtrelenmeyen, daha
doğrusu işin bu kısmını tamamiyle biz API kullanıcılarına bırakan kısa mesaj
gönderimlerinde karşılaşırız. Sık karşılaştığım ve internet üzerinde pek kaynak
bulamadığım için işin başa düştüğünü ve kendim .NET’in klasik Encoding’leri ile
çevirim uyumlu bir sınıf yazmam gerektiğini düşündüm. GsmEncoding aslında
Tasslehoff için yazdığım .NET kütüphanesinin bir parçası, kendisinin kaynak
kodlarına GitHub üzerinden erişebilirsiniz:

%[https://github.com/eser/tasslehoff-library/tree/master/Text]

Kodları kullanmak için .NET’in CLR’ında herhangi bir Encoding’i nasıl
kullanıyorsanız, yeni oluşturduğunuz bir `GsmEncoding` instance’ını da bu
şekilde kullanmanız yeterli. Örneğin:

%[https://gist.github.com/eser/5099144]
', 'tr          ', '01J6ZJ7PD8J8YP66QE91Q2DH5Q', '', '.NET ile GsmEncoding');
insert into "story_tx" ("content", "locale_code", "story_id", "summary", "title") values ('Konuyla ilgili twitter’da çok değerli Erhan Yakut’un tweet’ine yanıt vermemle
başlayan tartışmayı, söyleyeceklerim biraz uzun olduğu için buraya taşımak
istedim.

Tweet’i buraya almak gerekirse:

![](https://cdn.hashnode.com/res/hashnode/image/upload/v1659040552520/q8SgYR7DR.png)

%[https://twitter.com/yakuter/status/1134379969608540161]

En üstteki 10 milyon web sitesinden %34’ünün WordPress çalıştırdığı bilgisine
istinaden PHP’e ufak bir övgü var burada. Burada Internet’in en büyük 10 milyon
sitesinden bahsettiğimizin tekrar üzerine basayım, bu bilgiyi bir cebimize
atalım.

Ben ilgili tweet’e vermiş olduğum yanıtta şunu savundum:

%[https://twitter.com/eserozvataf/status/1134380938350202880]

Bunun karşılığında birçok kişi karşı arguman olarak “PHP’nin kolay
kurulabilirliği”, “5 TL’lik shared hosting’de çalışabilmesi” cümleleriyle
portatifliğini öne sürdüler. Şimdi bunlara yanıt vermek isterim;

Öncelikle PHP’nin portatifliği gerçekten takdire şayan. Bu portatifliğin
getirisi olarak artılarının/eksiklerinin olduğunu ve implementasyonunda
“uygulama mantığı” yer almadığı için birçok izolasyon sorunları yaşanabildiğini
şerh düşerek burada takdirimi yenileyeyim. FaaS seven ve evanjelistliğini yapan
biri olarak, PHP portatifliğini övmesem eksik kalır.

Ancak esas yanıtım şu: WordPress’in yaygınlığında elbette kolay kurulum ve ucuz
maliyet genel kullanım için bir etkendir; Ancak en üst 10 milyon web sitesinden
bahsederken bu %34’lük payda portatiflik etkenlerini gözlemleyebileceğimizi
düşünmüyorum. Türkiye’de dahi ortalama 2–3 dedike sunucu aracılığıyla host
edilen WordPress tabanlı web siteleri varken ve bu firmaların ceplerinden zaten
bakım maliyetleri çıkıyorken bu etkenlerin etkisi bana uzun süredir kurumsalda
çalışan bir profesyonel olarak gerçekçi gelmiyor.

Bir deneyimim ile devam edeyim, geçen yıl çalıştığım birime Türkiye’nin ilk
3’deki bir holding’in şirketlerini de etkileyecek çapta bir karar için uzman
görüşü soruldu. Buna istinaden CMS’ler ve platformlarla ilgili deneyimlerim
nedeniyle işi ben üstlenip bir rapor hazırladım. Bu iş öyle bir çaptaydı ki
holding bünyesinde “içerik” kelimesi geçtiği anda artık nihai olarak seçilen
çözüm kullanılıyor olacaktı. Hatta buna bağlı olarak bünyesinde olduğum firmanın
dışarı vermiş olduğu teklifler dahi bu altyapı gözönünde bulundurarak
verilecekti.

Rapor tamamlandığında yaklaşık 7–8 çözümü marketlerinden desteğine, lisanslama
politikalarından güvenlik açıklarına birçok alanda değerlendirmiştim. Vermiş
olduğum uzman görüşünde WordPress’i 5 üzerinden 4,5 puanlamama ve
değerlendirmede WordPress’e en çok yaklaşan ekosistemin puanı 2,5-3’lerde
kalmasına rağmen nihai karar WordPress olmadı.

Çünkü WordPress’in yumuşak karnı PHP idi. Ne bir holding gibi bir yapı için
hemen elde edilebilir, çok ciddi garantiler veren (kişisel garantilerin ötesinde
şeylerden bahsediyoruz) bir yetkinliğe ulaşabiliyorlardı ne de PHP’nin dolaylı
ya da dolaysız şekilde bir güvenlik açığı oluştuğunda/yük aldığında buna
karşılık verebilecek bir operasyon — kurumsal destek planları oluşturulabilirdi.

Raporumda “kolay kurulabilir” ve “portatif” gibi zaten operasyon işi olan
kısımları değerlendirmeye almadığımı belirtmeme gerek yok sanırım.

Sözün özü tekrar başa geleceğim. Internet’in en üst 10 milyon sitesinden
bahsediyoruz. Etkenler çok daha farklı.

Bu nedenle ben WordPress’in başarısını “başarılı bir ürün” olmasına bağlıyorum.
Aksi takdirde Ghost (herhangi bir rakibi) JavaScript ile değil PHP ile
yazılsaydı WordPress kadar başarılı olabilirdi diyebilirdik. Ancak böyle bir
yaklaşım yalnızca WordPress’in mevcut başarısını küçük görmek anlamına gelir.
', 'tr          ', '01J6ZJ7SSNZ0JT71YTCB1YJ8F2', 'Twitter''da başlayan tartışma üzerine fikirlerimi yazmak istedim; WordPress''in başarısını PHP''nin kolay kurulumu veya ucuz maliyetine değil, gerçekten iyi bir ürün olmasına bağlıyorum. Profesyonel deneyimlerime göre, büyük ölçekli sitelerde "portatiflik" yerine kurumsal destek ve güvenilirlik önem kazanıyor.', 'WordPress’in kerameti');
insert into "story_tx" ("content", "locale_code", "story_id", "summary", "title") values ('Haberlerden merhaba', 'tr          ', '01J3YPPRB2EWZK7N31DMC50Y8C', 'Bu yazı AYA''nın yeni özelliğinden bir selamlamadır.', 'Merhaba dünya!');
insert into "story_tx" ("content", "locale_code", "story_id", "summary", "title") values ('Çokça kez “yeni başlayan ve başlamak isteyenler” için içerik üretmediğim için
eleştirildiğim, hatta konferanslara uygun bulunmadığım (ileri seviye konular
seçiyormuşum) oldu. Ben her zaman kendim üretirken mutlu olduğum içerikler,
konulara odaklandım.

Bugün Murat Demirci ile konuşurken, bu konu hakkında yazmak istediğimi fark
ettim. Ve sanırım ilk kez “yeni başlayacak olanlara” naçizane görüşlerimi ifade
edebileceğim bir konuyu ele alacağım.

---

20’li yaşlarımın 8-9 yılı üniversitelerde geçti. Üniversiteye girmeden önce
alaylı olmuş, yazılım dillerine hakim biri olarak, bu süre içerisinde neredeyse
her yıl şöyle bir döngü yaşadım:

Yazılım geliştirme gereksinimi olan derslere katılan, ancak programlama
bilgileri sınırlı olan birçok arkadaşım da dönem başlangıçlarında “Eser ne olur
bu dönem çalışalım, ben öğrenmek istiyorum artık” diyerek kendilerince “benim
sponsorluğumu” almış olsalar da, ancak kritik sınavların önceki gününe kadar hiç
aramaz sormazlardı.

Sınavdan önceki gece ise ev arkadaşlarımın hoş görüsüne sığınarak evde 9-10 kişi
ders çalışırdık. Ben içimden iflah olmaz bir idealistlikle “neyse, en azından
midterm/vize’de iyi not alamayacak olsalar da finallere yetişirler” derken; bir
süre sonra tüm yarı dönemin müfredatını bir gecede odaklanıp çözmeye çalışan
arkadaşlarım sanki onlara ben zorla bir şeyler öğretmeye çalışıyormuşçasına bir
tavra bürünüp; neredeyse hep bir ağızdan “yaw şimdi aslında yarın hoca ne sorar
onlara kafa yorsak... o soruları ezberlesek olmuyor mu? bunları sonra da
öğreniriz” diyiverirlerdi.

Bildiğim kadarıyla hiçbiri hayatına “yazılım geliştirici” olarak devam etmiyor
bugün.

---

Bir yandan da benim hikayem şu şekilde gelişti:

Yazılım geliştirmeye bilgisayarla bir şeyler yapabilmek için başladım. Önce
BASIC, Pascal, C öğrendim. Sonra Visual Basic, Delphi. Photoshop, Corel Draw ve
Flash öğrenerek önce grafikerliğe soyundum. Sonra arkadaşlarımın ısrarıyla
HTML/CSS ile web sitesi yapmayı da portföyüme kattıktan sonra oyunlara modlar
yapmaya başladım. Oyun camiasından tanıştığım birisi (Khan-F) bana ASP öğretti,
JavaScript’i kendim öğrendim. Sonra daha yetenekli ve daha yaygın diye PHP’e
geçmem gerekti.

Para kazanmam lazım dedikten sonra ise PHP ile başladığım işler Python, C#
sonrasında DevOps pratikleri, ondan sonra yöneticilik vasıfları, ürün geliştirme
bakış açıları, Agility ve Software Craftsmanship derken bugünlere öğrene öğrene
ve proje yapa yapa geldim.

---

Herkesin öğrenme, öğrendiğini özümseme yöntemleri farklıdır. Yukarıda kendi
hikayemi yazmış olmam lütfen “herkes böyle yapsın” gibi algılanmasın. Geleceğim
nokta bu değil. Bu yazıda bahsedeceğim ve eleştireceğim nokta “öğrenmeyi bekleme
durumu” olacak.

Toplulukta ve yaptığımız sunumlarda “nasıl başlayabiliriz” sorusu çokça kez
herkese soruluyor. Bunun artık kemikleşmiş, hatta kendisinin komedisi haline
gelmiş bazı şablon yanıtları bile var: “İngilizce ile başla” gibi.

Ancak bu soruyu bugün bana, yarın başka birine, ertesi gün başka birine
soranların motivasyonu daha çok yukarıda bahsettiğim “yarınki sınavı geçecek
kadar konuya zaman ayırmak isteyen” arkadaşlarımla benzeşiyor. “Öyle bir kilit
yanıt alayım ki, hem çok basit olsun hem de yapınca yazılım öğrenmiş olayım”
gibi.

Maalesef bunun bir süreç olduğunu, gün gün gelişme kaydedileceğini, nasıl hiç
vücut geliştirmeyle ilgilenmemiş birinin bir günde muhteşem kaslı bir vücuda
sahip olamayacağı, eline hiç gitar almamış birinin bir günde herhangi bir
şarkıyı çalabilir hale gelemeyeceği gibi kabullenmek gerekiyor.

Kapak resmini “göklerden gelen İsa Peygamber” koyma nedenim de, “bir gün gelecek
kurtarıcıyı bekleme” davranışına yaptığım vurgu aslında.

Bugün “bekleme” konumunda olan hiçkimse, birden bir aydınlanma yaşamayacak ve
birden tüm yazılım külliyatını anlamıyorken anlayabilir hale gelmeyecek
maalesef. O nedenle efor gerektirmeden, işin “yarın birden her şeyi çözeceğiniz”
bir yanıtı varmış gibi bir çözümü varmışçasına “kolayını” aramayı bırakıp
öğrenmeye başlamak gerekiyor.

Öğrenmeye başlamak da biraz “ellerini kirletmek”ten geçiyor. Uzaktan izleyerek
değil, örneğin bir takım projeleri alıp bozsalar bile bazı denemelerde bulunan
arkadaşlar oluyor. Onları çok takdir ediyorum.

Örnek vereceğim, “[açık kaynak](https://acik-kaynak.org/)’ı açık kaynak
yapıyoruz” parolasıyla açık-kaynak.org’da arada
[GitHub üzerinde](https://github.com/acikkaynak/acikkaynak-website) bazı
maddeler açtık, birçok insan katkıda bulundu. Bazı noktalarda takılanlar da
Telegram’dan benimle iletişime geçip yardım istiyorlar. Elimden gelen yardımı
yapmaya çalışıyorum zamanım el verdiğince.

Projelerle başlamak isteyenler için kaynak bolca var, örneğin
[açık kaynak projeler](https://acik-kaynak.org/projects/) veya
[good first issue](https://goodfirstissue.dev/)’den bunlara da bakabilirsiniz.

Hackathonlara katılanlar, YouTube’dan bir şeyleri izlerken bir yandan
öğrendiklerini taklit edip yeni bir şeyler yapmaya çalışanlar, kitaplardan
okuyup kavramsal olarak konuları özümseyip sonra onu yazılıma dökmeye
çalışanlar, v.b. hep kendi yöntemlerince bir şeyleri “ellerini kirleterek”
yapmaya çalışıyor.

Ellerini kirletenler her gün birer adım atıp 3 ayda 5km yol yaparken,
ışınlanabileceğini düşünenler ve ışınlanmayı bekleyenler maalesef yol gidememiş
olacak.

O nedenle benim naçizane görüşüm, “hiç gelmeyecek ilham’ı” ya da “motivasyon
konuşmasını” beklemeyi bırakıp artık bir şeyleri öğrenme yolunda “ellerini
kirletmeye” başlamak olacak.
', 'tr          ', '01J6ZJ7NSFX726JYRV59WX72MP', 'Yazılım öğrenmenin aniden gerçekleşen bir mucize olmadığını; adım adım ilerleyen, pratikle güçlenen uzun bir süreç olduğunu anlatmaya çalıştım. Bu önkabulle birlikte nasıl hareket edilmesi gerektiğine dair de düşüncelerimi paylaştım.', 'Öğrenmek için “o anı” beklemek');
insert into "story_tx" ("content", "locale_code", "story_id", "summary", "title") values ('Son bir kaç haftamı .NET platformuna adayıp OData + Entity Framework ve Web API
ile geçirmişken “Data” kavramı üzerine düşünecek zamanım da oldu.

Bu sürede kazandığım deneyimlerimi sonra aktarmayı planlıyorum, ama öncesinde
“uzak servis üzerinde yer alan izole veri” isimli bulmacada eksik parçanın ortak
bir entity tanımlama formatı olduğuna kanaat getirmiş durumdayım.

Microsoft dünyasından uzaklaşmaksızın yorumlarsak; .NET Framework 5 ile Entity
Framework 7’ün gelişi yakın. EF7 daha hafif olma iddiası ile “Code First”
dışındaki tüm yaklaşımlardan bünyesini arındıracağını duyurdu. Bu kararı
elbetteki kullanıcıların class aracılığı ile entity tanımlamaya olumlu tepki
vermelerinden dolayı kolayca alabiliyorlar. Fakat halen bu alanda araç kutusu
eksiksiz görünen Microsoft’un dahi kendi “Entity Data Model” (EDM)si hem Entity
Framework hem de OData kütüphaneleri tarafından “paslaşır” bir biçimde
kullanılamıyor.

Resmi daha da genişletip, JavaScript yazan Frontend Developer’ından, bir API’ye
entegre olmayı bekleyen fakat PHP kullanan müşteriye kadar gittiğimizde hem
herhangi bir ORM üzerinde POCO sınıflarını dil-bağımsız oluşturabilecek,
veritabanına verdiğimizde anında fiziksel olarak tablo/görünüm/ilişkileri
oluşturabilecek, hatta Visio gibi bir tool’a verdiğimizde bize UML çizmesini
sağlayacak bir data tanımlama dili yok kabul görmüş olan.

Data yapılan işin bu kadar merkezinde iken, data’yı ve data ilişkilerini halen
her platforma özgü tanımlıyor olmak büyük sorun bana kalırsa.
', 'tr          ', '01J6ZJ7QSGP0K3EDZM0CNCKN74', '', 'Eksik parça — Data tanımlama');
insert into "story_tx" ("content", "locale_code", "story_id", "summary", "title") values (E'Ne yazık ki kapsamının dışına taşmış, ağızlara sakız olmuş bir kavramdır
“tekerleği yeniden keşfetmek”. Bir kaç pazarlama kaygısıyla yazılım
eğitimlerinde, seminerlerinde “catch phrase” yani reklam sloganı olmuş bu cümle
ile sürekli karşılaşan biri olarak bir inceleme yapmak isterim.

Öncelikle şu anda halk ağzındaki kullanımıyla ünlü “tekerleği yeniden keşfetme”
örnekleri ile başlamak isterim.

**Google Chrome**\\
Netscape’den sonra tekerleği yeniden keşfetmemiş Firefox tüm özgür
yazılımcıların gönül desteğini alarak herkesin popüler browser’ı oldu. Hem
extension desteği hem de standartlara uymasıyla insanlar Firefox fan’ı olmuştu
ki, Google şirket politikası gereğince kendi browser’ına sahip olmayı istedi.
“Tekerleği yeniden keşfetmiş” KDE ekibinin önce Nokia sonra Apple desteği ile
ilerleyen projesi “WebKit”, tekerleği tekrar keşfederek “Chromium” isimli yeni
bir arabirim ile şu anda pazar payı %80’lerde olan tarayıcı Google Chrome’un
tabanını oluşturacaktı.

**Linux**\\
UNIX sistemi varken ve çalışmaktayken 70lerin sonunda BSD çekirdeği mevcuttu.
Sun da 92 yılında Solaris’i yayınlarken ondan bir yıl önce halihazırda elinde
GNU’ya ait komut seti olan, 87’de yazılmış minix varken “tekerleği yeniden
keşfetmeye” istekli Linus Torvalds Linux’u yayınladı.

**Git**\\
CVS, SVN varken yine Linus Torvalds tekerleği yeniden keşfetme güdüsüne karşı
koyamayarak Git’i ortaya attı. Git şu anda en popüler açık kaynak
paylaşım/geliştirme sistemi GitHub’ın omurgasını oluşturmakta.

**Windows**\\
Apple cihazlarında kullanıcı arabirimiyle insanları mest ederken Microsoft
tekerleği yeniden keşfetmek isteyerek Windows’u yayınlamaya başladı. Yine aynı
Windows, kendi işletim sistemi ailesinin stabil olmadığını düşünerek Windows NT
adını verdiği sistemi kendi tekerleğini yeniden keşfederek geliştirdi ve Windows
XP’den itibaren yeni nesil Windowslarda bu tabanla devam etme kararı aldı.

**MySQL**\\
PostgreSQL’den tutun, Oracle’a milyonlarca veritabanı varken MySQL çıktı. Üstüne
üstlük en büyük özelliklerinden biri yeniden tekerlek keşfetmek isteyenler için
veritabanı motoru geliştirmeye açık olmasıydı.

Daha kaç örnek verebilirim bilmiyorum. En yaygın örnekler üzerinden gitmeye
çalıştım. Yukarıdaki projelerin neden “tekerleği yeniden keşfetmek” olmadığını
anlamak için önce keskin bir ayrım olan “ürün ve teknoloji farkı”nı özümsememiz
gerekecek.

Teknoloji, benim kendi ezberden tanımım ile “bilimin ışığı altında insanların
hayatlarını kolaylaştıran geliştirmeler”dir. Öyle cep telefonlarına gitmemize
gerek yok, oturduğumuz sandalye insani bir problemi çözdüğü için bir
teknolojidir. İçinde fizik kurallarını gözetecek kadar bilim de vardır,
insanlara yarar sağlayacak kadar geliştirme de. Ve evet, “tekerlek” bir
teknolojidir, hatta özdeyişleri bile onunla kuruyorsak en popüleridir belki de.

“Tekerleği yeniden keşfetmek” bilimin geldiği, test ettiği, uyguladığı “kesin”
bir süreci tekrar yaşamamayı adresler. Bu bir kural değildir bu arada, şayet bu
bir kural olsaydı eminim üniversite dediğimiz kavram göçer ve insanlar “tez”
dediğimiz çalışmaları yapıyor olmazlardı. Özdeğişin burada bize bahsettiği bazı
kavram ve varsayımları tekrar sorguya açmamaktır. Sandalye yapıyorsak,
sandalyeyi kullanan kişinin bir oturak kısmı üzerinde durması gerektiği artık
bilimseldir. Deney edilmiş, denenmiştir.

Eğer ağzına “tekerleği yeniden keşfetmek”i sakız edenlerle mevcut
sandalyelerimizi üretenler aynı fikre sahip olsaydı ilk üretilen sandalye
modelini halen kullanıyor olurduk. Yani el/kol koymak için kenarları olan,
istediğimiz yöne döndürebildiğimiz, sırtımızı yaslayabildiğimiz, hareket
ettirebildiğimiz, hava alan sandalyeler olmazdı. Çünkü teknoloji ile ürün/model
arasındaki farkı özümsemekle ilgili bir sorun var ortada.

Soyut örnekten biraz daha bilişim sektörüne geri dönersek, Google Chrome ortaya
çıkarken HTTP protokolünü yeniden keşfetmedi, istediğimiz internet adresini
yazdığımız adres çubuğunu da. Ama bunları alıp tekrar kendi deneyimleriyle
yorumladı. Daha hızlı olacağını inandığı için HTTP yerine SPDY’yi önerdi. Adres
çubuğunu arama çubuğu ile birleştirerek Omnibox adı altında hem adres hem arama
yapabileceğimiz bir forma dönüştürdü. Çünkü Google’ın bir vizyonu, yazılım
geliştirmeyle ilgili bir akımı ve bir geleneği/deneyimi mevcuttu. Tüm bu
bilgiler ışığında kendi ürününü veya browser modelini ortaya attı. Keza Linux de
öyle, Git de, MySQL de…

Şu anda kendi projelerini yapan insanları “tekerleği yeniden keşfetmek”le
kınamaya yakın bir şekilde yargılayanlar, genellikle çalıştıkları yazılım
şirketlerinde daha önce başka birileri için yazılmış yazılımların 100. kopyasını
yazmaktalar. Sistem zaten yıllardır içerik yönetim sistemi, e-ticaret altyapısı
vesaire üretiyor ve yeni yeni tekerlekler keşfediyoruz. Hadi startup olduğunu
varsayalım, yeniden Facebook yazmaya gerek yok derken bile Facebook’a ortada
Yonja varken üye olmuş neslin bireyleri olarak çelişkiye düşüyoruz.

Çözümüm bilim dünyasını ele almak olacak, kullanıcı deneyimi, iyi yazılım,
yazılım geliştirme metodları gibi konular bilim gibi sürekli geliştirilebilir,
sürekli sorgulanabilir kavramlardır. Eğer bilimde mevcut fikirleri tezlerle
sorgulayabiliyor, alternatif üretebiliyorsak ilerleme kaydetmek için tekerlek
başta olmak üzere bir çok şeyi yeniden keşfetmemiz gerekecektir. Koca koca
firmalar işin ucunda maddi kaygılar varken bile projelerini çöpe atıp fikirleri
kopyalarken, tekerleği yeniden keşfediyor olmuyor ama bu duruma çok takılan bir
kesim mevcut.

Sözü ironi ile bitirme ihtiyacı duydum, tekerleği yeniden keşfetmemeli ama
piyasada bir çok lastik markası mevcut.
', 'tr          ', '01J6ZJ7R7GBWBFNK9ME20DJEVE', '', 'Tekerleği yeniden keşfetmek');
insert into "story_tx" ("content", "locale_code", "story_id", "summary", "title") values (E'PHP’nin 5.1 sürümüyle birlikte “Standard PHP Library” (SPL) isimli bir kavram
tanıtıldı. SPL içerisinde, PHP geliştiricilerinin çok kullanıldığına kanaat
getirdikleri bir takım hazır tip ve fonksiyon tanımları mevcut. Bu tanımlar PHP
bünyesinde servis edilerek, bir çok kütüphane ve yazılımcının bir standart
üzerinden çalışmasına imkan veriyor. Ne yazık ki bir çok geliştirici SPL’den
habersiz olarak geliştirme yapmaya devam ediyor.

[PHP Manual’ın SPL ile ilgili bölümü](http://tr.php.net/spl)nü paste etmektense
burada bir özet geçmek isterim.

### Veri tipleri

[**SplStack**](http://tr.php.net/manual/en/class.splstack.php): Bir array
üzerinde _array\\_pop_, _array\\_push_ dışında bir işlem yapmıyorsanız array
alternatifidir.

[**SplQueue**](http://tr.php.net/manual/en/class.splqueue.php): Bir array
üzerinde `array_shift`, `array_push` dışında bir işlem yapmıyorsanız array
alternatifidir.

[**SplPriorityQueue**](http://tr.php.net/manual/en/class.splpriorityqueue.php):
Bir array’e eklediğiniz elemanların önemlilik sıraları varsa ve bu elemanları
array’i okurken her zaman önem sırasına göre almak istiyorsanız array
alternatifidir.

[**SplFixedArray**](http://tr.php.net/manual/en/class.splfixedarray.php): Bir
array’e ekleyeceğiniz eleman sayısını önceden biliyorsanız array alternatifidir.

### Iteratorler

[**DirectoryIterator**](http://tr.php.net/manual/en/class.directoryiterator.php):
Bir dizin içeriğini okur.

[**LimitIterator**](http://tr.php.net/manual/en/class.limititerator.php): Başka
bir iterator’un belirli eleman aralığını kullanır. (İlk 5 eleman gibi)

### Interfaceler

[**Countable**](http://tr.php.net/manual/en/class.countable.php): Sınıfın
`count` fonksiyonu ile sayılabilir olmasını sağlar.

[**Serializable**](http://tr.php.net/manual/en/class.serializable.php)\\*:
Sınıfın `serialize` ve `unserialize` fonksiyonlarıyla çalışabilmesini sağlar.

[**ArrayAccess**](http://tr.php.net/manual/en/class.arrayaccess.php)\\*: Sınıfın
array gibi davranmasını, oluşan nesnenin `$nesne[‘eleman’]` biçiminde erişimini
sağlar.

[**Traversable**](http://tr.php.net/manual/en/class.traversable.php)\\*: Sınıfın
foreach ile dolaşılabileceği bilgisidir. `Iterator` veya `IteratorAggregate` ile
birlikte kullanılır.

[**Iterator**](http://tr.php.net/manual/en/class.iterator.php)\\*: Sınıfın
`foreach` ile dolaşılabilmesini sağlar.

[**IteratorAggregate**](http://tr.php.net/manual/en/class.iteratoraggregate.php):
Yine sınıfın dolaşılabilmesini sağlayan bir iterator sınıfına sahip olduğunu
belirtir.

### Design Pattern yardımcıları

[**SplObserver**](http://tr.php.net/manual/en/class.splobserver.php): Observer
design pattern’deki gözetleyen ve güncellemeleri alan sınıf interface’i.

[**SplSubject**](http://tr.php.net/manual/en/class.splsubject.php): Observer
design pattern’de gözeticilerin abone oldukları sınıf interface’i.

### Fonksiyonlar

[**iterator\\_to\\_array**](http://tr.php.net/manual/en/function.iterator-to-array.php):
Elinizdeki bir iterator’ı foreach’le dönmeksizin array’e dönüştürür.

[**iterator\\_apply**](http://tr.php.net/manual/en/function.iterator-apply.php):
Bir iterator’ın her elemanı için parametre olarak verdiğiniz fonksiyonu
çalıştırır.

[**iterator\\_count**](http://tr.php.net/manual/en/function.iterator-count.php):
`count` ile sayamadığınız bir iterator’ın eleman sayısını verir.

### Exceptionlar

Kendi Exceptionlarınızı yalnızca `Exception` sınıfından değil, bir çok yer
sınıftan türetebilir böylece istisnai durumları yakalarken daha net tip
ayrımlarında bulunabilirsiniz. Hatta benim önerim, PHP gibi interpreter tabanlı,
runtime’ı request-response ile sınırlı (tek tüketimlik) bir ortam için mümkün
olduğunca mevcut sınıfları kullanın. SPL’nin önemi de bu noktada belirginleşiyor
zaten. Exceptionların anlamlarını teker teker yazmak isterdim ama isimleri zaten
yeterince açıklayıcı, burada önemli olan bu Exceptionların inheritance yapısı, o
da şu şekilde:

[Exception](http://tr.php.net/manual/en/class.exception.php)

Son olarak, kesin bir çizgi olmamakla birlikte, ben kendi yorumumca:
`LogicException`‘ı kullanıcı/programcı tarafından yanlış bir işlem yapıldığında,
`RuntimeException`‘ı ise bir işlem gerçekleştirirken bir zorluk oluştuğunda
kullanıyorum. `ErrorException` zaten PHP’nin hata oluştuğunda kullandığı özel
bir sınıf.
', 'tr          ', '01J6ZJ7QH103QMEQZK0QAQWJWY', '', 'PHP SPL — PHP Standard Library');
insert into "story_tx" ("content", "locale_code", "story_id", "summary", "title") values ('Bildiğiniz üzere Agile’ı benimseyen ürün geliştirme ekipleri bir kaizen (sürekli
iyileştirme) kültürü içerisinde _retrospective_ toplantıları yaparlar, bunu da
sonraki döneme ilişkin bir _planlama_ ritüeli takip eder. Yerelleştirip,
“takkeyi önüne alıp düşünmek” olarak da ifade edebileceğimiz bu “değerlendirme”
ve “planlama” pratiklerini yararlı bulduğumdan, periyodik olarak kendi
gönüllülük esasında yaptığım işlere de uyguluyorum. Amaç sürekli iyileştirmek ve
yeni bir şeyler denemek olunca bu kez bu değerlendirmeyi yazılı ve şeffaf
gerçekleştirmek istedim.

Bir başka deyişle bu yazıdaki her şey aslında kendi özelimde olan “kişisel”
değerlendirmelerim. Bu yazı kapsamında yalnızca neyi, nasıl ve neden yaptığımı
netleştirmek için yazılı hale getirmeye çalışacağım. Belki bugün net yanıtlar
alamayacağım ancak bu ilk adımdan başlayarak, kişisel dönem değerlendirme
sürecimin sonunda nihai olarak şu sorulara yanıt almaya çalışacağım;

- Hangi alanlarda gönüllü faaliyet gösteriyorum?
- Bu alanlardaki çalışmalarımın kapsamlarını daha da netleştirebilir, konsantre
  hale getirebilir miyim?
- Bu alanların hangilerinde gerçekten yararlı oluyorum?
- Bu alanların hangilerinde geri dönüş alıp kendimi geliştirme fırsatı
  buluyorum?
- Bu alanlardaki çalışmalarım bir bütüne, bir ilkeye, bir fikre hizmet ediyor
  mu?
- Bu alanlardaki çalışmalarım ile bütünsel bir ilke, bir yaklaşım oluşturabilir
  miyim?
- Bu alanlardaki çalışmalarımın birinden vazgeçerek daha yararlı olabileceğim
  yeni bir alan bulabilir miyim?
- Değerlendirmelerimin sonucunda bir yol haritası çıkartabildim mi? Bir hareket
  planı oluşturabilir miyim?

### Konferans/Meetup Konuşmaları

Bu aslında kendimi oldukça sınırlamaya çalıştığım bir alan. Çocukken konuşma
zorluğu yaşayan biri olarak konuşmalar yapmanın beni ekstra motive eden bir yanı
da var. Ancak benim bu alandaki faaliyetlerimin diğer gönüllü olduğum
kısımlardan zaman çalmaması, ve ürettiğim yararın sürdürebilirliğinin korunması
gerekiyor. Bu nedenle;

Bu konuşmaları sayısal anlamda senede en fazla 5 etkinlikle sınırlıyorum.

Aynı başlıktaki konularla ilgili farklı konuşmacılar oluşmaya başladığında ben
bırakıyorum. Bir örnek olarak Serverless/FaaS konusunu 2 yıl önce yalnızca
[Daron Yöndem](https://daron.me) ve ben anlatırken, şimdi çok daha fazla kimse
konu hakkında konuşma gerçekleştiriyor. Dolayısıyla ben daha yararlı olacağım
farklı konular oluşmadıkça kendimi tekrar etmemeye çalışıyorum.

Yine aynı kitleye seslenen aynı şehirdeki organizasyonlara aynı konu ile
gitmemeye çalışıyorum. Bir örnek olarak Frontend İstanbul’da TypeScript
anlattıysam ve jstanbul’da TypeScript’i tekrar anlatma şansım olursa; ilk olarak
bakacağım şey katılımcıların aynı olup olmadığı olacaktır. Aynı kitleye birden
fazla kez aynı şeyi aktarmanın bir yararı bulunmuyor. Bu nedenle üniversitelerin
lokal öğrenci etkinliklerine daha fazla açığım bu dönemde.

Dolayısıyla insanlarla tanışmak, etkinliklerin heyecanlarını paylaşmak ve
sürdürülebilirliklerini desteklemek için 2019’da da konuşma davetlerine yanıt
vereceğim ancak bunun için YouTube’de etkin olacağım zamandan çalmak
istemiyorum.

**Hareket Planı:** Belirlediğim sınırlamalar dahilinde etkinliklerde konuşmalara
devam edebilirim.

### YouTube Yayıncılığı

Geçtiğimiz sene başladığım, yatırım yaptığım, ufaktan kitle oluşturmaya
başladığım ve daha da büyüteceğim gönüllü çalışma alanım aslında YouTube.

Ben başlarken daha ıssızdı, ancak canlı yayın streamleri halen yazılım içeriği
anlamında potansiyeline ulaşamadı. Twitch gibi platformda yazılımla ilgili
içerik aradığınızda “Creative” başlığı altında live coding dışında bir içerik
bulamıyorsunuz. Ancak önümüzdeki dönemde bunun değişeceğine inanıyorum. Kendimi
de bu değişimin öncülerinden biri olarak konumlandırmaya çalışıyorum. Bu nedenle
YouTube benim için çok kıymetli bir gönüllülük alanı.

YouTube için hareket planımı çıkarttım aslında. Canlı yayınlarda zaten oyun
temelli twitch yayıncılığının pratikleri literatüre kazandırılmış durumda. Bunu
yazılım için de kullanmak mümkün. Bilgisayarlarımızın başında kahvelerimizi
yudumlarken,

- yazılım ürünlerini unbox edelim,
- bir radyo/podcast yayınıymışçasına yazılım gündemini konuşalım,
- açık kaynak proje ve araçları keşfedelim,
- birlikte konferans kayıtlarını izleyelim,
- ve bir araya gelip networking oluşturalım,

istiyorum. [YouTube Kanalım](eser.live) üzerinden takip etmeye başlayıp siz de
katılımcı olmaya başlayabilirsiniz.

**Hareket Planı:** YouTube üzerine yatırıma devam.

### Evangelist Yaklaşım

Gerek yüz yüze, gerek internet üzerinden iletişimde bulunduğum yazılımla ilgili
herkese bir takım ilkelerin ve yaklaşımların evangelistliğini yapmaya
çalışıyorum. Ancak direkt “kendi önceden hazırlanmış görüşümü” aktarmak yerine,
insanları önce dinleyip daha sonra farklı bir yol/yöntem mümkün mü birlikte
konuşmak burada benim icra ettiğim en önemli farklılık.

Örneğin daha çok kurumsal iş ortamlarında bulunuyorsam, o disiplinde kendine
henüz pek de yer edememiş ancak yavaştan nüfuz etmeye başlayan Açık Kaynak
pratiklerinin savunuculuğunu yapıyorum. Gerçekten organizasyonlarımızı güvenlik
duvarlarının, kapalı kod tabanlarının arkasına mı saklamalıyız?

Yine insanların ne yaptıklarını çok iyi bildiği, özverili olduğu ancak toplamda
verimden bahsedilmediği ortamda “ne üzerine”den çok “nasıl çalıştıklarını”
sorgulayan Agile pratiklerini ortaya atmak başkalarına bir bakış açısı
kazandırmak anlamında yararlı olabiliyor.

Biraz daha geleneksel yazılım yaklaşımlarını gördüğümde, unix felsefesinin,
DevOps yaklaşımlarının, cloud computing’in iddiaları ile bu görüşleri sınamaya
çalışıyorum. Bazı konularda insanlarla yaptığınız bir konuşma o an verimli
olmasa dahi bir tohum olarak bir fikrin filizlenmesine, daha sonra da büyümesine
neden olabiliyor.

İnsanlarla konuşmak, dertlerini ve içinde bulunduğu durumları dinlemek, onlara
yeni bir şeyler önermek kesinlikle beni mutlu ediyor. Temelde mühendisliğin
çözümcü yaklaşımını sosyal konular üzerinde de ele almak hem bana bu alanda bir
sorun-çözücü olarak deneyim katıyor hem de insanlar pek kendi görüşlerini
dinleyen birilerini bulamadığından bir boşluğu dolduruyor.

**Hareket Planı:** Devam.

### Açık Kaynak Projelerim

Benim için [acik-kaynak.org](http://acik-kaynak.org)‘un geliştirilmesi başlı
başına bir “açık kaynak” proje oldu, bu projenin sahibi ben değilim fakat
inisiyatifi başlatan bendim. Bu noktada konuyla ilgili iki çuvaldızım var, biri
kendime. Benim için acik-kaynak.org‘un geliştirilmesi hep ertelendi. Halen
gerekli zamanı bulup geliştirmeye devam etmek istiyorum ancak bunu
başaramıyorum. Diğer yandan ülkemizde dillendirildiği kadar açık kaynak
sevdalısı olmadığını iletilen proje sayılarından görmüş bulunuyorum.

[Tassle](https://github.com/eser/tassle) benim uzun yıllar üzerinde çalıştığım
bir .NET Standard component paketi, .NET Core’un çıkışı ile birlikte ciddi
refactorler geçirdi. Ancak halen istediğim olgunluğa kavuşamadı.

[Ponyfills](https://github.com/eser/ponyfills),
[Immunity](https://github.com/eser/immunity),
[Evangelist](https://github.com/eser/evangelist),
[Enthusiast](https://github.com/eser/enthusiast),
[ServiceManager](https://github.com/eser/servicemanager),
[CoFounder](https://github.com/eser/cofounder),
[Senior](https://github.com/eser/senior),
[Maester](https://github.com/eser/maester),
[Consultant](https://github.com/eser/consultant),
[JsMake](https://github.com/eser/jsmake) gibi ufak tefek JavaScript toolchain
olarak ifade edebileceğim projelerim üzerinde vakit buldukça çalışıyorum.

Diğer yandan [Darty](https://github.com/eser/darty),
[Create-Darty-App](https://github.com/eser/create-darty-app),
[Darty React App](https://github.com/eser/darty-react-app),
[Darty React Hooks App](https://github.com/eser/darty-react-hooks-app) gibi kod
tabanı oluşturma aracı projem devam etmekte. Aslında ts-spa-boilerplate isimli
frontend geliştirme boilerplate’imin create-react-app pratiklerinin çok
benimsenmesi sonrasında evrimleştiği nokta olarak ifade edebiliriz “Darty”yi.

Tabii ki bunlar aktif maintain ettiğim ve son dönemde ilgilendiklerim. Daha
sayamadığım ve GitHub üzerinde bulunan 103 adet açık kaynak repository’m
bulunuyor, bunlar arasında oyun geliştirme alanında sayılabilecek
[Quake 3 Mod](https://github.com/eser/q3now)’undan tutun, C++’da yazdığım AST
üzerinde çalışan parserlara veya Solidity üzerinde oluşturduğum Etherium
SmartContract’lara kadar birçok farklı proje bulunuyor. Tek kişi bunları
maalesef çok etkin maintain edemiyorum. Bazılarını seçip bu projelerin
geliştirmelerine daha konsantre devam etmem gerekecek.

**Hareket Planı:** Üzülerek bir sınırlama getirmek, daha konsantre gitmek
gerekiyor.

### Akademik Çalışmalar

Bu yazıyı yazarken Kuzey Kıbrıs’ta blockchain ile ilgili lisansüstü tezimi
savunmak ve teslim etmek üzere bulunuyorum. 36 yaşında biri olarak okul ve
eğitim her zaman bırakmadığım bir parça oldu hayatımda. Tatil dönemlerinde bile
edX üzerinden MIT’den ders aldığım, hem çalışıp hem de eğitimimi sürdürdüğüm
dönemler hayatımın bir gerçeği gibiydi.

Şimdi ise son mezuniyetimle birlikte, özel hayatıma daha fazla zaman ayırmak ve
kendime iş ve kültür alanlarında bir şeyler katmak için bu aktif eğitim kısmını
bırakmak zorunda hissediyorum.

Oluşacak entelektüel kazanım boşluğunu belki daha fazla kitap ve yayın tüketerek
telafi edebilirim.

**Hareket Planı:** Farklı alanlara zaman ayrılacak.

Aslında her başlığın altında kendi değerlendirmem ve hükmüm yazdığı için yazının
sonuç kısmında yazacak pek bir şey kalmıyor.
', 'tr          ', '01J6ZJ7NGMK53TAMW0C783FV0B', '', 'Gönüllü olduğum işlere dair planlama ve değerlendirme');
insert into "story_tx" ("content", "locale_code", "story_id", "summary", "title") values ('Gelecekle ilgili vizyon tanımları hep ilgimi çekmiştir. Hatta geriye dönük bir
tarama yapsak blog yazılarımın satır aralarından bununla ilgili bir çok ipucu
çıkacaktır. Bunlara bir yenisini ekleyeyim: çocukken televizyonda X-Men’den
sonra en sıkı takip ettiğim seri Beyond 2000 idi (pazar günleri TRT’de).

Gelecekle ilgili ön-okumalarımdan ve iddialarımdan biri şu ki; Gelecekte
bilgisayar, bilgi ve iletişim sistemlerinin kırılamaz, kandırılamaz olup
olamayacağını hayatımızı bir kaç yönde etkileyecek. Şifreleme ve güvenlik
sistemlerinin geleceği aynı zamanda bizim de geleceğimizi pratik yönden
şekillendirecek.

Halihazırda her gün eğitim, finans, telefon hizmetleri, üretim, reklamcılık
alanlarında bilgisayarın rolünün giderek arttığına şahit oluyoruz. Ama henüz
“medya”sı hariç sosyal ve bireysel hiçbir kavram bilgisayarize edilmemiş
durumda.

Bu “bilgisayarize” edilmemiş kavramlar ve imkanlar üzerine bir çok kurgu
oluşturulabilir. Ama herkesin demokrasi üzerine daha fazla kafa yormaya
başladığı bugünlerde, bu konuların aralarından “bireysel temsil”i ve dolayısıyla
demokrasiyi konu olarak belirlemek istiyorum. Konuyu genişletmek adına
“güvenilirliğini tarafların gözleme özgürlüğüne dayandıran klasik demokrasi
yerini bilgisayar-destekli modern bir demokrasi modeline bırakabilir mi?”
sorusuna yanıt vermeye çalışarak başlayacağım.

Hatta organize gidelim, önce sorunu belirleyeceğim ardından vizyon tanımını
yapacağım, fizibilitesine ise siz karar verirsiniz.

**Sorunlar:**

Bence, demokrasi halkın iradesinin temsil edildiği bir sistem olarak; günümüzün
artan nüfüs ve toplumdaki farklı görüş sayısı ile klasik yöntemlerle üstesinden
kalkamayacağı bir dinamizme ihtiyaç duyuyor.

Özellikle gelişmekte olan ülkelerde, hızlanan yaşam koşullarınca yasama organına
düşen yük artıyor. Hemen uygulanması beklenilen ve hızla meclisden geçirilen
yasaların “sarsıntı” yaşatmaması, kendi içerisinde nispeten tutarlı olması adına
istikrar gereksinimi duyuluyor. Bu da zaten “öcü” görünen “koalisyon”
hükümetlerinden uzak, oy verenlerin gidişata 4–5 senelik periodlarla müdehale
edebilecekleri, temsil oranının yaklaşık %40 ila %50’sinin dışarıda kaldığı
parlementoları oluşturuyor. Ayrıca bireyler tarafından seçilen
milletvekillerince yürüyen, ama nüfus kalabalıklığından ötürü bir kere dahi
seçilen ve seçenin aynı ortamda bulunulmasına ve bu kararların tartışılmasına
izin vermeyen parlementer sistem temsil gereksinimini sağlıklı karşılayamıyor.

Aynı sistem “oy verenin görüşünün veya ihtiyacının değişmesine neden olan”
hiçbir olaya refleks oluşturamıyor. Güven oylamalarının öncesinde ve sonrasında
ekonomik krizler ülkelerdeki istikrarlı gidişatı bozuyor.

**IT Çözümü:**

Mükemmel bir analiz yapmamış olsam da, IT uzmanı olarak çözümüm e-devlet
entegrasyonuna bağlı çok ileri bir güvenlik doğrulaması ile kullanıcıların
erişebildikleri “özel” sayfalarda milletvekili, belediye seçimlerini
yapabilmesi, istedikleri zaman bunları değiştirebilmeleri. Ve yine yapılan
icraatlere ait feedbackler ve referandumlar bu sistemin içerisinde yer alabilir
diye düşünmekteyim.

Bu tarz bir sistem aslında politikacıların mevcut şartlarda sağlık
politikalarına mı yoksa işçi haklarına mı daha fazla önem vermesinin
beklendiğini onlara iletmesi, gerçekleşen icraatlerinin de nasıl tepkiler
doğurduğunu takip edebilmesi için hükümet/muhalefet yönlü de fayda
sağlamaktadır.

Bir nevi “bilgisayar ve temsil yardımlı açık demokrasi”den bahsettiğimin
farkındayım. Bu durumun bazı sakıncaları yok mu? Tonla yeni problem var!

En başında da ülke ekonomisine olan etkileri. Anlık tepkilerin istikrar
bozmasını engellemek için “kamu oyu” biraz “geç” veya “periodik” olarak
yansıtılabilir diye düşünüyorum.

Geliştirilebilir fikirler birer başlangıç noktasına gereksinim duyarlar. En
azından tartışmaya veya hayal etmeye açarak bugünden Bilgisayarlı Demokrasi’yi
planlayabileceğimizi düşünüyorum.
', 'tr          ', '01J6ZJ7RG7RCMAFHWRX01BCDEM', '', 'Bilgisayarlı Demokrasi');
insert into "story_tx" ("content", "locale_code", "story_id", "summary", "title") values ('İyi bir kelime seçimi gibi duyulmasa da “yarı vizyon” tam olarak anlatmaya
çalışacağım senaryoyla oturuyor.

Microsoft’un firma ve yaptıkları işler hacminde büyüklüğü tartışılmaz. Fakat
birçok başarısız girişimini gözlemliyor, ciddi oranda profesyonelin geçim
kaynağını sağlayan ekosistemlere sahip olmalarına rağmen düşük güvenilirliğe
sahip olduğunu biliyoruz. Bu durum firmanın iyi yaptığı işlerin yanısıra bir o
kadar kötü yaptığı işlerin de olduğu konusunda bize ipucu veriyor.

Yazı aslında Web Teknolojilerine odaklı fakat Microsoft’un bulunduğu konumla
ilgili çarpıcı bir örnek paylaşmak istiyorum, mutlaka hikayenin bir kısmını siz
de duymuşsunuzdur. Steve Jobs henüz iPad ile çıkagelmemişken Microsoft’un
2000lerin başında “Tablet PC”yi insanlara teşvik etme çabaları vardı. Yıl
2001’di, şu anda televizyonda reklamlarını son bir yıldır izlediğiniz
katlanılabilir notebook’larda olduğu gibi tablet pc’nin ekranı katlanılabiliyor,
henüz dokunmatik ekranların yaygın olmadığı o tarihte bir kalem yardımıyla
Windows’u veya ARM işlemcili kopyaları için Microsoft’un Windows CE ismiyle
geliştirdiği mobil işletim sistemini (iOS ve Android’den önce!) sunuyordu
sizlere. 2001’de Microsoft’un bu arzına talep gelseydi olaylar çok farklı
gelişecekti. (Kaynak:
[Wikipedia](https://en.wikipedia.org/wiki/Microsoft_Tablet_PC))

Gelelim bu güne, benim gibi Web Teknolojileriyle de yakınen ilgilenen
yazılımcılar şu anda fırında ısınmayı ve gelecek senelerde karşımıza çıkmayı
bekleyen yeni web standartları için heyecanlanıyordur.

Angular 1.5, Polymer, React gibi kelimeleri duyduysanız Web Components;
Electron’u, SPA’yı duyduysanız HTML Application kavramlarını biliyorsunuzdur. Bu
iki kavram da Microsoft’un henüz JavaScript’i ECMAScript standartlarında kabul
etmeyip JScript olarak isimlendirdiği, “DHTML” ismi çevresinde web uygulamaları
konusunda kullanıcı tabanını heyecanlandırmaya çalıştığı, Internet Explorer’ın
en büyük pazar payına sahip olduğu dönemlerde düşünülmüş ve Internet Explorer’a
hatta bir kısmı Windows’a entegre edilmiş konseptler.

HTA ([HTML Application](https://en.wikipedia.org/wiki/HTML_Application)):
Çalıştırılabilir bir HTA dosyası aracılığı ile tıpkı masaüstünüzde çalışan diğer
uygulamalar gibi web uygulamanızın bir pencerede çalışması sağlanıyor. Üstelik
farklı alt sayfalara linkler de yok, Microsoft diyor ki tek bir HTML ve script
yığını ile tek HTML sayfadan oluşan bir uygulama yapabilirsin. Ekstra olarak
bazı sistem çağrıları ve pencere özelliklerine ulaşabiliyorsunuz.

HTC ([HTML Components](https://en.wikipedia.org/wiki/HTML_Components)):
Oluşturduğunuz bir element’e (veya html tag diyelim) bir davranış tanımlayıp
JavaScript aracılığıyla nesneye gelen olay/eventleri o bileşene özel
değerlendirebiliyorsunuz. Angular’ın 1.5, React’ın ve Polymer’in yapmayı
sağlandığı şey de bu değil miydi?

Şu soruyu sormadan edemiyorum: Internet Explorer 5.5 gibi 2000’de çıkmış bir
browser’ın getirdiği yenilikleri Microsoft topluluğa benimsetseydi ne olurdu? Bu
kavramların özelliği veya benim anlatmaya çalıştığım “Microsoft ilk yaptı”
değil, bu konseptlerin şu anda W3C standartlarına girmiş, yani topluluğun
istediği ve benimsediği kavramlar olması. 15–16 sene önce DHTML üzerine
çalışmalarım olmasa bugün ben de bu konseptlerin yeni yeni ortaya atılıyor
olduğunu düşünecektim ama değiller.

Buradan hareketle Microsoft’un “ileriyi” görüp ihtiyaçlara göre bir çözüm
ürettiğini fakat zamanlama, motivasyonunu insanlarla paylaştırma, insanlara
fikrini benimsetme gibi konularda oldukça eksik kaldığını gözlemleyebileceğimizi
düşünüyorum.

Buraya kadar “vizyonerlik” tanımı bozulmuyor esasında, “yarı vizyon” tespitim
bundan sonrası ile ilgili. Microsoft da ilk harekete geçtiği andan itibaren
kavramı geliştirmek yerine vazgeçmeye yakın bir çizgide duruyor. Silverlight,
Windows Phone, WinFS gibi Microsoft’un yalnızca son 3 senede vazgeçtiği
ürünlerden ciddi bir liste oluşturabilmek mümkün.

Artık topluluğu çağıran değil, topluluğun olduğu yere gitmeye çalışan daha
akıllı bir Microsoft var gibi görünüyor ama bunun henüz TypeScript dışında bir
ürünü olduğunu gözlemlemiyorum. Umarım zaman bize Microsoft’un potansiyelini
gerçekleştirebildiği daha çok örnek gösterir.
', 'tr          ', '01J6ZJ7MJDAG3GZHY49SE509RD', '', 'Microsoft’un yarı vizyonu üzerinden WebComponents ve SPA');
insert into "story_tx" ("content", "locale_code", "story_id", "summary", "title") values (E'Microsoft kullanıcı deneyimlerinde 2010’a kadar grafik arabirimi kullanımını
arayüzünde herhangi bir ayrıma gitmeyi tercih etmedi. Hatta bu süreç içerisinde
mümkün oldukça her işlemin grafik arabirimi üzerinden gerçekleştirilebilmesi
yönünde inisiyatifler aldı.

.NET Core ve Windows 10 ile birlikte bu kararlılıkta ciddi değişimler
gözlemleyebiliyoruz. Örneğin Windows 10 Anniversary Update 1 ile birlikte
bash’ın Windows’a gelişi, Windows içerisinde tanımlanan bir “Developer Mode”
seçeneği, .NET Core’u yüklemek için dahi komut satırı araçlarının önplanda
tanıtılması v.s.

Önümüzdeki süreçte Microsoft kendi teknolojilerine odaklanmış yazılımcıların da
Visual Studio dışında bir yazılım geliştirme ekosisteminde, diğer teknoloji ve
parçalarla uyum içerisinde çalışmasının önemli bir gereksinim olduğunun farkında
olsa gerek.

Bu araçların büyük bir çoğunluğunun çalıştığı platform tabii ki, macOS, Linux
gibi işletim sistemlerinde olduğu gibi komut satırı olacak. Ve nasıl
masaüstünüzü, start menüsünü rahat kullanmak günlük bir ihtiyaç ise komut
satırında da rahat hissedebildiğiniz kadar yazılım geliştiricileri kendilerini
üretken hissedecekler.

Ben bu noktada kendi sıfırdan kurulum reçetemi paylaşıp, windows komut satırında
hangi araçları kullandığımı listelemek istiyorum.

**1\\. Chocolatey**

[https://chocolatey.org/install](https://chocolatey.org/install)

Paket yöneticisi program kurulum ve güncelleme yönetimi için kritik önem
taşıyor. Apt, Yum, Brew gibi örnekler zaten işletim sistemlerinin önemli bir
parçası iken son dönemlerde Microsoft’un da Chocolatey projesine desteği ile
Windows için 3. parti geliştirilen bu aracın şu sıralar oldukça popülerlik
kazandığından bahsedebiliriz.

Ben aşağıdaki diğer programların kurulumunu chocolatey üzerinden yapıyorum.

**2\\. Clink**

[https://chocolatey.org/packages/clink](https://chocolatey.org/packages/clink)

bash kullanıcıdan gelen komut girişlerini okumak için _readline_ isimli bir
kütüphane kullanır. Clink basitçe anlatmak gerekirse aynı deneyimi yaşamanız
için _readline_’ı windows’un komut satırı arabirimini sağlayan _cmd.exe_’ye
uyarlar.

**3\\. Sudo, Make, OpenSSH, Curl, Vim, Git**

[https://chocolatey.org/packages](https://chocolatey.org/packages)

Sırasıyla,

Komutları Administrator hakları ile başlatabilmeniz için _sudo_,\\
 Makefile’ları windows ortamında çalıştırabilmeniz için _make_,\\
 Komut satırından putty gibi yazılımlara gerek kalmaksızın ssh bağlantısı
kurabilmek için _win32-openssh_,\\
 Komut satırından http istekleri yapabilmek için _curl_,\\
 Komut satırında gelişkin bir editor olan _vim_,\\
 Yine GitHub for Windows veya herhangi bir IDE/text editor desteğine gerek
duymaksızın Git erişimi için _git_,

paketlerini chocolatey üzerinden indirebilirsiniz.

**4\\. DOSKey ile aliaslar ve diğer konfigurasyonlar**

[https://gist.github.com/eser/e8eb6f8e7f57418a1775b2b68960a0a1](https://gist.github.com/eser/e8eb6f8e7f57418a1775b2b68960a0a1)

MS-DOS’un 5.0 sürümüyle hayatımıza giren DOSKey aslında komut satırı için bir
makro aracı. Halen günümüzde de Windows komut satırında yer aldığı için
kendisini bash’da yer alan `alias` komutunun bir benzeri işlevle kullanıyorum.
Bu sayede evde kullandığım macOS’dan işyerinde kullandığım Windows’a geçtiğim
anda komut satırına `ls` yazdıktan sonra hatamı `dir` yazarak düzeltmek zorunda
kalmıyor, sistemler arasındaki farkı minimize edebiliyorum.

Mevcut kullandığım makrolar aşağıdaki gibi:

```
DOSKEY ls=dir /B $* DOSKEY ll=dir $* DOSKEY clear=cls DOSKEY open=start $* DOSKEY service=net $2 $1 DOSKEY apt-get=choco $*
```

Maalesef `cmd.exe` kendi çalıştığı anda belirli başlı komutları bizim için
işletmiyor (`.bashrc` veya `bash_profile` benzeri bir işlevden bahsediyorum). Bu
nedenle yukarıdaki makroları açılan her komut satırı için tekrardan tanımlamam
gerekiyor. Eskiden aşina olduğumuz `AUTOEXEC.BAT` çoktan tarihe karışmış
durumda. Bu nedenle ben kendi kullanıcı dizinime bir `env.cmd` dosyası
oluşturarak her komut satırını açtığımda elle kendisini çalıştırıyorum.
Yukarıdaki linkten benim kullandığım `env.cmd`nin içerine ulaşabilirsiniz.

Diğer bir değişiklik ise benim artık sıkıldığım, belki sizin de değiştirirek
keyif alabileceğiniz prompt. Sizden komut isteyen `C:\\Windows\\system32>`
yazısını değiştirmek istiyorsanız,
[Scott Hanselman’ın A better PROMPT for CMD.EXE or Cool Prompt Environment Variables and a nice transparent multi-prompt yazısı](http://www.hanselman.com/blog/ABetterPROMPTForCMDEXEOrCoolPromptEnvironmentVariablesAndANiceTransparentMultiprompt.aspx)nı
takip ederek gerekli değişiklikleri gerçekleştirebilirsiniz.

Benim tercihim aşağıdaki gibi:

```
SET PROMPT=$+[$m$p]
```
', 'tr          ', '01J6ZJ7PPXDJMG7H8M8K2C1BTF', '', 'Windows Komut Satırı için kullanılabilecek güçlü araçlar');
insert into "story_tx" ("content", "locale_code", "story_id", "summary", "title") values ('Bana göre işim hiçbir zaman bir çözüm sağlanana dek kod yazmak olmadı. Hatta
maaşımın karşılığını geliştirme yaptığım sürelerde değil, bunun dışında kalan
sürelerde hakkettiğimi düşünüyorum. Buradaki argümanım çok basit: “Hedef çözüm
sağlamak ise karmaşık veya sorunlu bir işi nihayete erdirebilmek için karmaşık
bir çözüm ortaya koymak için çalışabilirsiniz. Veya önce problemi sadeleştirip
onun dengi sadelikte çözüm üzerinde çalışabilirsiniz.”.

Karmaşıklıktan uzaklaştırılan ama yönetilmeye devam edilebilen her iş zamanımız
başta olmak üzere her türlü sınırlı sermayemizden edilen net kazanç benim için.

[Twitter](https://twitter.com/eserozvataf)’da, herhangi bir
[konuşma](../../../presentations/tr/README.md)da denk geldiyseniz genellikle
Serverless, Cloud ve DevOps çözümleri konusunda az da olsa evangelist
yaklaşımımı fark etmişsinizdir. Bu konuların beni heyecanlandırmasının nedeni,
hepsinin temelinde işi gerçekleyenler için karmaşık olan problemleri onlardan
uzaklaştırıp, daha soyut bir üst katmanda çözümler üretmelerine imkan veriyor
olmaları. Bugün tam da bu kavramların kesişimiyle ilgili bir şeyler aktarmak
istiyorum. Bir yazılımın teslimat süreci, kulağımıza daha aşina gelecek
ifadelerle belirtmek gerekirse “derleme işleminden başlayarak canlıya çıkana
kadar gelişen olaylar”.

2006’da çalıştığım devBiz’e girerken şirket ortaklarından ve benim için değerli
isimlerden biri olan Hakan Eskici bana Mike Gunderloy’un “Coder to Developer:
Tools and Strategies for Delivering Your Software” kitabını vermişti. İlk kez
yazılım için entegrasyon (integration), teslimat (delivery) ve yükleme
(deployment) tanımları ile orada karşılaşmıştım.

Mike Gunderloy — Coder To Developer

Bu kavramları ezberimden aktarmam gerekirse,

**Integration**: Bir geliştirme sürecinde yapılan tüm ekleme ve değişikliklerin
yazılımın tamamına ve bağlılıklarına olan uygunluğu/entegrasyonunu sınayan ve
sağlayan süreçtir. Genellikle çıktısı test edilmiş, derlenmiş yazılım dosyaları
olur.

**Delivery**: Bir yazılımın paketlenmesi ve bu paketin hedefine ulaştırılma
sürecidir. Genellikle çıktısı lisansları, konfigürasyonu, imzaları ve sürüm
bilgisi eklenmiş halde kuruluma hazır yazılım paketidir. Docker Image’dan tutun,
bir zip dosyasına herhangi bir biçimde paketlenmiş olabilir.

**Deployment**: Bir yazılımın bir sisteme/ortama yüklenmesi ve yüklendiği
ortamda çalışmaya başlayarak şekilde ayağa kaldırılması sürecidir.

Siz bu kavramları son yıllarda başlarına bir “Continuous” ifadesi ile birlikte
duymaktasınız, **Continuous Integration**, **Continuous Delivery** ve
**Continuous Deployment**. Bu da yapmış olduğunuz bir değişikliği merkezi bir
Git/SCM sunucusuna gönderdiğiniz anda bu süreçlerin otomatik olarak başlaması ve
tanımlanan görevlerin sizin tetiklemenize gerek kalmaksızın kendiliğinden
işletilmesi, eğer bir hata ile karşılaşılırsa da bundan haberdar olmanız
anlamına geliyor.

Bu süreçlerin tümünü işler hale getirmek ve işletmek günümüzde DevOps
yetkinliğine ve/veya kimliğine sahip kimselerin yükümlülüğü. Ancak her ne kadar
her gün bir yenisi tanıtılan DevOps araçları ve pratikleri ile DevOps
süreçlerinde hedeflere ulaşılsa da; elinizdekilerle bir süreci uçtan uca
tasarladığınızda mutlaka aklınızdakini/beklentinizi birebir uygulayamadığınız
birçok durum oluşuyor. Hatta elinizdeki senaryo ve hizmet bağımlılıkları
nedeniyle çok temiz, esnek ve olan biteni kavramanın zaman almayacağı bir yapı
yerine kompleksiteye boğulmuş olmaktan kaçılamıyor.

Konuyla ilgili Fatih Kadir Akın’la bir konuşmamızdan sonra durumun vehametini
şöyle bir tweet’le paylaşmıştı kendisi:

%[https://twitter.com/fkadev/status/1074972412184354816]

Peki neyi sorunlu görüyorum ve ne iyileştirilebilir? Hemen açıklayayım…

Geliştiricinin ve DevOps’un işi paket hazırlayana kadar olabilirdi. Elimizde
geliştiricilerin müdehale ettiği bir kod tabanı var. Bu kod tabanının
derlenmesi, test edilmesi, buradan docker image oluşturma v.b. paketlenme
işlemlerinin _pipeline as code_ pratiğiyle DevOps’cular tarafından da
sağlandığını düşünelim. Buraya kadar bir sorun yok.

Ancak bu noktadan sonra DevOps’cunun oluşan çıktının nereye, hangi senaryoyla,
ne zaman teslim edileceği ve yükleneceğini tasarlaması; buna bağlı olarak
başarısızlık, felaket ve sistem yetersizliği durumlarında sistemlerin nasıl
davranacağına dair de önlemler alması bekleniyor.

İşte bu noktada aslında DevOps idealize edildiği gibi hem yazılımın içerisini ve
işleyişini bilen, aynı yazılım ekibinden biri/birileri olamıyor. Çünkü yukarıda
bahsi geçen süreçler için AWS hizmetlerinden, Azure’a, Kubernetes’den Packer’a,
Spinnaker’dan Jenkins’e birçok araca hakim olmak zorunda bırakılıyor. Bu da bir
DevOps uzmanı için mesainin aslında ilgili yazılımın kod tabanından uzak, bu
araçlara daha yakın geçmesine neden oluyor.

Bu sürece uçtan uca baktığımda “paketleme” olarak ifade ettiğim noktadan
sonrasını tamamiyle “karmaşa” ve “üstlenilmemesi/devredilmesi gereken bir yük”
olarak görüyorum. Bir örnekle anlatmam gerekirse; esas işi e-ticaret olan bir
web sitesinin satın alınan ürünü paketledikten sonra bir araca binip müşterisine
kapıda teslimat yapmasına benzetiyorum bu durumu. Bu da beni bu yazının
başlığındaki teslimat’ı biz mi yapmalıyız sorusuna getiriyor.

Diğer yandan iyimser kalmamı sağlayan bir etken var: Cloud Hizmetleri.

Yukarıda Serverless’ın ismini geçirmemin nedeni temelde size sunduğu en büyük
hizmetin sizi uygulama mimarisi tasarlamaktan bir nebze kurtarıyor olmasıydı.
İleride AWS, Azure, v.b. servislerin birçok senaryoda yaptığı gibi teslimat
işlerinden de bizi kurtaracak çözümler sunacağına emin gibiyim.

Zaten bugün de yine anahtar teslim bir ürün söz konusu olmasa da Container
Registry ve Function hizmetleri kullanılarak bir app service’e deployment
yaparak benzer bir süreç işletebiliyorsunuz. Hatta Azure’un App Service for
Containers’ıyla birlikte Azure Container Registry’den webhook’lar aracılığıyla
Continuous Deployment yapmanıza izin veriyor. Belki daha farklı ürünler de
vardır ancak konumuz tam anlamıyla ürünlerin getirdiği çözümler değil.

Benim tartışmaya açtığım konu daha çok süreçteki rol ve sorumluluklar üzerine.
Ben derim ki; Developer ve DevOps’un sorumluluğu oluşturduğu paketi bir
container/file storage’a yüklediği anda sona ermeli. Geri kalan noktada
yazılımın ne kadar kaynakla, ne zaman, blue/green’le mi yoksa başka bir
stratejiyle mi yükleneceği tamamiyle hosting’i gerçekleştirdiğimiz servisin
orkestrasyonuyla ilerlemeli.

Bu uzun konunun ertesinde belki aşina olmadığınız birçok kavramdan bahsetmiş
olabilirim. Eğer bu konular ilginizi çektiyse, aklınızda soru işaretleri
oluştuysa onları dinlemek adına [bir YouTube yayınıma](https://eser.live)
beklerim İlerleyen günlerde özellikle pipeline oluşturmak, CI/CD entegrasyonu,
cloud hizmetleri, serverless gibi konularla ilgili içeriklere de yer vermeye
çalışacağım.

Bir sonraki yazıya dek, sağlıcakla kalın.
', 'tr          ', '01J6ZJ7R34NJXYD7X854XCH5AG', '', 'Teslimat bizim işimiz (mi?)');
insert into "story_tx" ("content", "locale_code", "story_id", "summary", "title") values ('[Fatih Acet](http://fatihacet.com/)’in Twitter üzerinden başlattığı “Uzaktan
Çalışma Deneyimi Olanlar, Ulaşın Bana!” çağrısının bir sonucu olarak bu yazıyı
blogluyorum.

Konu “Uzaktan Çalışma” olunca, aklımdaki her şeyi bir araya nasıl toparlayacağım
hakkında hiçbir fikrim olmamasına rağmen ilk kazmayı ben atmak istedim.

Halen birçok işveren açısından “ütopik”, yazılımcı açısından “sahip olunamayacak
kadar güzel bir hayal” gibi duyulan “Uzaktan Çalışma” modelinin sabit bir
reçetesi maalesef ki bulunmuyor.

Bu nedenle şirketler kendi kültürleri, olmazsa olmazları, işleri işe alış
yordamları ile “lokasyon bağımsız çalışanlar” gibi konseptleri sentezleyerek
kendi modellerini oluşturuyorlar. Bu modellerde oluşan bazı gedikler bazen
iletişimsizlik ile sonuçlanıyor, bazen ise kişilerin kendi özverileri ile
kapatılıyor.

Ben burada önce kendi deneyimlerimi işverenler ve iş arkadaşlarını konunun
merkezine alarak kronolojik olarak anlatacak, ardından da kendi çıkarttığım
sonuçları paylaşacağım. Yazının devamı olarak da daha kişisel organizasyona
yönelik farklı bir perspektiften konuyu değerlendireceğim.

#### a) Aynı Şehirde Uzaktan Çalışma

2008’in başlarında 15 ay askerlik sonrasında en özlediğim şeyin yazılım
geliştirme olduğuna karar verip zaman kaybetmeden iş aramaya başlamıştım. Askere
gitmeden önce çalıştığım şirketi (devBiz) Microsoft’un satın almış olması bana
hem biraz burukluk hem de yüksek özgüven veriyor, şartlar konusunda elimin güçlü
olmasını sağlıyordu. İyi maaş, Senior Developer pozisyonu, Cumartesi-Pazar
günleri çalışmama olmazsa olmazlarımdı.

Çalışmaya başladığım şirket mülakat esnasında “Cumartesi” talebime karşılık “45
saatini doldurduğun sürece bizim için hangi günler veya nereden çalıştığın
önemli değil” yanıtını vermişti. Bir iş merkezi içerisinde bulunan işyerine saat
8’den sonra fiziksel olarak erişemiyorduk, bu nedenle evden de çalışmalara devam
etmek için bu formulü bulmuşlardı. Sonradan takım lideri olup kaldırttığım,
“Timeless” isimli bir iş/zaman takip yazılımı kullanılmakta bu şekilde de iş
takibini gerçekleştirmektelerdi.

#### b) Farklı Ülkede Uzaktan Çalışma

Benim gibi 2002’de mezun olmuş her meslek liselinin üniversite kapısında
karşılaştığı sorun olan “katsayı problemi”ni yaşamış biri olarak, önüme koyulan
önlisans eğitiminin bana yetersiz geldiğini düşünüyordum. Özel bir üniversitede
eğitimime devam etme şansı ise 2010’un sonlarında benim için mümkün olmaya
başladı. Zaten ofis dışında çalışmalarımla takım liderliğine yükseldiğim firmaya
daha da uzaktan, başka bir ülkeden çalışmaya devam edecektim artık.

İletişimi canlı tutmak için tamamiyle yazılı ortama geçiş yapmıştık artık. O gün
olmayan Slack yerine konular kendi başlıklarında kalsın amacıyla bir forum
kullanıyor, toplantıları Skype üzerinden yürütüyor, şirketin sağladığı iPhone’la
anlık pushlanan e-mailler iletişimimizin bel kemiğini oluşturuyordu. Fakat
iletişim kalitesi beklentilerimi karşılamıyordu. Yüz yüze rahat iletişim
kurduğum, ofis ortamında bir şeyler paylaştığım insanların zamanla attığım
mesajları yalnızca savuşturmak adına yanıtlar verdiğini, kurulan cümlelerin
tonlarını çok fazla irdelediklerini hissetmeye başladığımda tüm bunlardan “o
gün” çıkardığım sonuç sorunun uzaktan çalışmakla ilgili olduğuydu.

#### c) Müşterilerle Uzaktan Çalışma

5 senelik bir süreye yaklaşacakken uzaktan çalışmakta olduğum firma ile
ilişkilerim koptu. Bulunduğum ülkedeki ilişki ağım Türkiye’dekinden daha az
olduğu için “Freelance” iş talepleri almaktaydım. Yüksek Lisans’a geçtiğim sene
benim için zorunlu bir yoldu bu artık.

Bazen aracı olan bir yüklenici, bazense direkt müşteri ile çalıştım. Buradan
çıkarttığım sonuç işin kendisi kadar iletişimin de hassasiyetle yönetilmesi
gerektiği, iki tarafın da beklentilerinin “ne kadar detaylı olarak
belirlendiği”nin de bu iletişimde kilit rol oynadığı oldu.

Daha önce
[Hayal kırıklığına uğramamak için sorulması gerekenler](./20160221-hayal-kirikligina-ugramamak-icin-sorgulanmasi-gerekenler)
başlıklı yazımda tarafların kendilerine sorması gereken soruların biraz
kenarından geçmiştim. Konuyu fazla da dağıtmadan ne kastettiğimi açıklamak için
beklentiler üzerinden müşteri ile ilişkimizi ele alalım:

- Müşteriniz verdiğiniz hizmetin kapsamını biliyor mu, bunu izah edebildiniz mi?
  (Sizin hatanız)
- Siz net olarak anlaşmada yapılacak işin kapsamını belirtmenize rağmen
  müşteriniz bunu okumayarak, ayağı yere basmamaya devam mı ediyor? (Müşterinin
  hatası)

Siz müşterinize “size talep ettiğiniz doğrultuda doğru içerik sağladığı takdirde
Google arama indeksinde yükselmeyi” vaad etmiş olabilirsiniz, müşteriniz bunu
anlamaksızın size para vermesinin yegane sebebini “Google’da 1. sırada çıkmak”
olarak belirlemiş olabilir.

Zamanla müşterileri kategorize etmeyi ve gerekirse kabul etmemeyi öğrenmek bu
tarz bir çalışma modelinde önem arz ediyor.

Yazının devamını
[Uzaktan Çalışma Deneyimim  —  Kişisel Organizasyon](./20161012-uzaktan-calisma-02)
başlığı ile bulabilirsiniz.
', 'tr          ', '01J6ZJ7MPJHY28R4M7KWF28SJE', '', 'Uzaktan Çalışma Deneyimim — İşverenler ve İş Arkadaşları (1/2)');
insert into "story_tx" ("content", "locale_code", "story_id", "summary", "title") values ('Sektöre yabancı olmayan herkesin bildiği fakat ara sıra sınırlarını belirtmek
adına tanımını düzgün yapmak ve yapılan tanıma referans verilmesi gereken bir
konu “açık kaynak”.

Tanımlara girmeden hemen şunu söylemem gerekiyor ki, “açık kaynak”, “özgür
yazılım” ve “kaynağı açık” kavramlarının anlamları ve kullanışları üzerinde
ortak bir uzlaşıya ulaşılabilmiş değil. Birçok kaynak bu ifadeleri eş anlamlı
kullanırken, işin detaylarında bu iki kavramından birbirinden ayrıştığını
savunanlar da mevcut.

Ben arada farklılıkların olduğunu düşünen gruba giren biri olarak, hem bu fikir
ayrılığının nedenlerini de elimden geldiğince açıklamaya hem de diğer
yazılarımda da olduğu gibi daha çok kendi bakış açımla bu kavramları yorumlamaya
çalışacağım.

### Kaynak

Sektörümüzde en çok kullanım alanı yazılımlar olsa dahi, herhangi bir çalışmayı
tekrar üretebilmemizi sağlayan her türlü materyale (döküman, veri, v.s.) kaynak
ismi verebiliriz. Microsoft Word dökümanından, bir programın kaynak koduna kadar
birçok şeyi kaynak olarak nitelendirebiliriz.

### Kaynağı Açık

Bir çalışmanın kaynağı kısmen veya tümüyle “herkes tarafından” erişilebilir ve
incelenebilir durumda ise o ürün için “kaynağı açılmış” diyebiliriz.

### Açık Kaynak

Her açık kaynak çalışmanın kaynağı aynı zamanda açıktır. Fark kelimelerin
yerlerini değiştirince değil,
[“açık kaynak”’ın Açık Kaynak İnisiyatifi tarafından yapılmış tanımı](https://opensource.org/definition/)nda
yatıyor.

“Kaynağı Açık” tanımına ek olarak, çalışmanın kaynağı “herkes tarafından”
değiştirilebilir ve tekrar dağıtılabilir durumda ise onun “açık kaynak”
olduğundan bahsedebiliriz.

Açık kaynak geliştirme ve kontrol gibi süreçlere “kaynak” olabilecek herkesi
dahil etmeyi hedefler. Bu sayede ortaya çıkan çalışma başka bir çalışmanın bir
parçası olabilir, birileri tarafından incelenerek hataları tespit edilebilir
veya yine birileri tarafından geliştirilerek daha da olgun bir hal alabilir.

### Özgür Yazılım

Her özgür yazılım aynı zamanda açık kaynaktır. Fakat özgür yazılım bundan çok
daha fazlasıdır. Çünkü özgür yazılımlar sahipliği reddeden bir politik anlayışın
ürünüdür. Dolayısıyla herhangi birinin veya birilerinin kontrolü altında
bulunmayan özgür yazılımların üzerinde yapılan değişiklikler de mutlaka herkesin
erişebileceği şekilde yayınlanmalı ve özgür yazılım olmayan başka bir programın
(hatta GNU Public License v3’e göre donanımın da) işlevsel bir parçası olarak
kullanılmamalıdır.

Özgür yazılımın olmazsa olmazları daha doğrusu “özgürlükleri” şöyle
alıntılayabiliriz:

- Herhangi bir amaç için yazılımı çalıştırma özgürlüğü (0 numaralı özgürlük).
- Her ne istiyorsanız onu yaptırmak için programın nasıl çalıştığını ögrenmek ve
  onu değiştirme özgürlüğü (1 numaralı özgürlük). Yazılımın kaynak koduna
  ulaşmak, bu iş için önkoşuldur.
- Kopyaları dağıtma özgürlüğü. Böylece komşunuza yardım edebilirsiniz (2
  numaralı özgürlük).
- Tüm toplumun yarar sağlayabileceği şekilde programı geliştirme ve
  geliştirdiklerinizi (ve genel olarak değiştirilmiş sürümlerini) yayınlama
  özgürlüğü (3 numaralı özgürlük). Kaynak koduna erişmek, bunun için bir
  önkoşuldur.

Kaynak: [GNU](https://www.gnu.org/philosophy/free-sw.tr.html)

Konuyla ilgili daha fazla detay için, hayatını özgür yazılıma adamış
[Richard Stallman’ın açık kaynak ve özgür yazılımın farklılıklarına değindiği yazısı](http://www.gnu.org/philosophy/open-source-misses-the-point.en.html)nı
okuyarak detaylara inebilirsiniz.

Son olarak konu özgür yazılım olduğunda, şu motto akıllarda yer etmeli: “free as
in speech, not as in beer”.

### Lisanslar

Peki bir çalışmayı ne “kaynağı açık”, “açık kaynak” veya “özgür yazılım” yapar?
Cevap “lisans” ismi verdiğimiz, kaynağın sahibi tarafından oraya koyulmuş,
çalışmayı eline geçirene bir takım izinler veren veya sınırlar çizen bir belge.

Eğer çalışmanın sahibi hukuksal olarak duruma çok hakim biri olarak kendi özel
şartlarını koymadı ise karşılaşacağınız lisanslar çoğunlukla MIT, BSD, Apache,
GPL v.b. lisansları gibi herkesçe bilinen hazır bir şablonda olacaktır.

### Açık Kaynak neyi hedefler?

İşin son olarak felsefesinden bahsetmek gerekirse, açık kaynak şeffaflığın ve
paylaşımcılığın merkeziyetçilik/kapalılıktan çok daha yararlı ve verimli olduğu
düşüncesini savunur.

Açık kaynak “bilim” gibi ilerler, alanında uzman kişilerin ortaya attığı
çalışmalar çürütülerek/düzeltilerek/desteklenerek mükemmeli arama çabasında bir
veri oluştururlar.

Şu anda günlük kullandığımız veya yararlandığımız birçok yazılımın açık kaynak
dünyası tarafından geliştirdiğini göz önünde bulundurarak açık kaynağın başarılı
olduğunu söyleyebileceğimizi düşünüyorum.
', 'tr          ', '01J6ZJ7N3G5HWRNDZ8YRJD8CDJ', '', 'Açık Kaynak Felsefesi');
insert into "story_tx" ("content", "locale_code", "story_id", "summary", "title") values (E'“Futurologist” Richard Watson 50 yıl içerisinde kullanımı sona erecek şeylerin
bir listesini oluşturmuş. Her ne kadar kendimi internet gazetesi haberciliği
olarak nitelendirdiğim bu tarz “catchy” şeyleri paylaşırken hayal edemiyor olsam
da, 90 jenerasyonu mensubu olarak “geleceğe ve bilime olan merak” nedeniyle bunu
paylaşmak istiyorum.

> – 2009: Tamir.\\
> – 2012: Tüplü televizyon.\\
> – 2013: Faks makineleri.\\
> – 2014: Kaybolmak.\\
> – 2016: Emeklilik.\\
> – 2017: Cep telefonu, mendil.\\
> – 2018: Sekreter, DVD.\\
> – 2019: Postane, Kütüphane, Kasap, Bedava Araç Parkı.\\
> – 2020: Patentler.\\
> – 2022: Bloglar.\\
> – 2024: Haftasonu Tatilleri.\\
> – 2025: Ücretsiz Otoyollar, AM Radyo.\\
> – 2026: Cilt Kırışıklıkları.\\
> – 2027: Öğle Yemeği, Gizlilik Hakkı.\\
> – 2028: Birleşimler.\\
> – 2030: Anahtarlar, Reality Şovlar.\\
> – 2033: Bozuk Paralar.\\
> – 2035: İşlenmemiş Yağ, Spam E-Mailler.\\
> – 2036: Bağımlılıklar.\\
> – 2039: Sağırlık, Karbon Emisyonu, Ulusal Para Birimleri, Kağıt Para,
> Cüzdanlar, Uyku.\\
> – 2040: Kravat.\\
> – 2045: Oturma Odası.\\
> – 2049: Gazeteler, Körlük.\\
> – 2050’den Ötesi: Ölüm, Acı, Çirkinlik, Plastik Cerrahi, Ulusal Devletler.
', 'tr          ', '01J6ZJ7NXRZC49DCQFFK8RV16R', '', '50 yıl içinde tükenecekler');
insert into "story_tx" ("content", "locale_code", "story_id", "summary", "title") values ('3 Ocak 2015’te
[Microsoft’a güçlü bir komut satırı için açık mektup](./20150103-microsofta-guclu-bir-komut-satiri-icin-acik-mektup)’u
yazdığımda açık kaynak organizasyonların yazılım geliştirme alanındaki
yenilikleri geliştirmede başı çektiği bir dönemi yaşıyorduk.

[Açık Kaynak Felsefesi](./20160330-acik-kaynak-felsefesi) ile üretilen yazılım
geliştirme ortamları ve araçları daha fazla kullanıcı tarafından benimsenmeyi,
kendi kendini ilerletebilen topluluklar oluşturabilmeyi, daha hatasız ve direkt
olmayı başarabiliyordu. İş yazılımcının tercihine kaldığında Microsoft ve
benzeri kapalı kaynak kuruluşları ancak maddi bir kazanım varsa bir alternatif
haline geliyordu.

Fakat her ikisi birden mümkündü. Apple başta olmak üzere hem açık kaynak
araçları kendi ekosistemine dahil etmiş, yazılımcılara tanıdık ve iyi test
edilmiş araçlar sunan firmalar da vardı. Microsoft da bunlara dahil olmaya
çalışıyordu.

Microsoft o dönem Visual Studio’nun “hafif siklet” sürümü Visual Studio Code’u
duyurmuş, .NET’i Linux/OS X üzerinde çalışacak şekilde .NET Core ismiyle yeniden
tasarlayan ve bunların tümünü yazılımcıların ayağına giderek GitHub üzerinden
gerçekleştirmeye başlamıştı.

Bunun üzerine yukarıda bahsettiğim yazımda Microsoft’un daha da ileri giderek
GNU araçları ve Bash’i Windows’a getirmesinin gerekliliğinden ve PowerShell’in
yetersizliğinden bahsettim.

Yazımı tamamlayıp paylaştığım mecralarda aldığım tepkiler “Microsoft gibi büyük
firmalardan direkt paraya dönüştüremeyeceği değişimlerin beklenemeyeceği”
üzerinde yoğunlaşmıştı.

Evet, Microsoft’un daha çok kurumsal ve kitlesel çözümler üreten bir firma
olduğunu aşikar, dolayısıyla “paraya dönüştürme” argümanının altında da haklı
sebepler yatıyor.

Fakat bir noktayı da kaçırmamak lazım. Microsoft aynı zamanda yazılımcıları
elinde tutmayı başararak bugünlerine geldi. Kendi ekosistemi için “yazılımcı
yetiştiren” ve yazılım tarafında da bazı “çöpe gitmesi istenmeyecek kadar büyük
markalar”a sahip bir firmadan bahsediyoruz.

Her yazılımcı kurumsal bir dünyada uyanmıyor mesleki kariyerine, lisansı olmayan
bir IDE ile, belki evinde belki okulda belki bir ajansta, HTML/CSS, PHP, Python,
Ruby v.s. üzerine çalıştıktan sonra Microsoft dünyasına giriş yapan yazılımcının
Microsoft için oluşturduğu izlenim oldukça önemli. Son zamanlarda Microsoft’un
içerisinde dahi Macbook’u ile sunumlara katılan yazılımcıların altını çizdiği
gerçek Microsoft’un “orada duran”ı kendi ekosistemini kullananlara
sağlayamadığını ve tercih edilmeyen olduğunu görüyoruz.

Anlaşılan o ki Microsoft benim gibi yukarıdaki mektubu yazacak kadar “isyan
eden” yazılımcılara duyarsız kalmadı. Ve bugün Windows 10’un bir sonraki
güncellemesi ile, yaratıcı üretimin çoğunu sağlayan yazılımcıların neredeyse yok
saydığı CMD veya PowerShell’i değil olması gerektiği gibi bash’i
kullanabileceğimizi açıkladı.

Özetlemek gerekirse Microsoft gelecek dönemde bu açılımla da sınırlı kalmayıp;
Ubuntu’nun yayımcısı Canonical ile bir ortaklığa başlayacak ve Linux üzerinde
çalışan Wine’ın tersi olarak nitelendirilebilecek şekilde Windows’a “Linux
Subsystem” isimli Linux binarylerin de Windows’da yorumlanmasını sağlayabileceği
bir yapı getirecek.

Detaylar elbette ki devam eden BUILD 2016 ile daha da ortaya çıkacak.

Ancak konuya geri dönüp toparlamak gerekirse; İstatistiksel olarak GitHub’da
olan ve kurulum/operasyon aşamasında komut gereksinimi duyan projelerin bash’e
ve diğer GNU araçlarına göre dökümante edildiğini bilmesek de (eminim çok büyük
bir orandır) Microsoft dünyanın geri kalanına yaklaşmak için büyük bir adım attı
diye düşünüyorum.

[BUILD 2016’daki ilgili açıklamanın videosu](https://youtu.be/GgutlKV42SY?t=4434):
%[https://youtu.be/GgutlKV42SY?t=4434]
', 'tr          ', '01J6ZJ7KSEQ6ERWW48T3BZCQRY', '', 'Microsoft’un güçlü komut satırı — Perde 2');
insert into "story_tx" ("content", "locale_code", "story_id", "summary", "title") values ('Daha önce
[Cloud computing’de son kullanıcı açısından gizlilik](./20111128-cloud-computingde-son-kullanici-acisindan-gizlilik)
başlıklı yazımda gizlilik üzerine olan fikirlerimden biraz bahsetmiştim. 2011’de
yazdığım yazı üzerinden geçen zaman içerisinde teknoloji kullanıcıların bu
konudaki duyarsızlığı nedeniyle malesef ki daha aç gözlü bir biçimde
verilerimizi kullanır hale geldi.

Şu anda internette yayınladığınız içerik için yasal olarak suçlanabilirsiniz.
Bunu tartışmıyorum, benim tartıştığım bu içeriğin aslında size ait olmaması.
Yani bir şeyin sorumluluğu bizde ama aynı zamanda o şey üzerinde kontrol sahibi
değiliz. Internette bulunan sayfalar sizin ürettiğiniz, size ait olan bir şeyi
saklayıp size vermeme yetkisine sahip.

Bir alışveriş merkezine girdiğinizi düşünün, yemek yediniz, sinema izlediniz ve
günü tamamladığınızı düşünüp eve döneceksiniz. Bu noktada güvenliğin binadan
ayrılmanıza izin vermediğini, buna gerekçe olarak da girişte sağda yazan “giren
bir daha çıkamaz” tabelasını gösterdiğini düşünün.

Şu an yasalar böyle bir alıkoymayı engelliyor, çünkü bedeniniz sizin ve
hukuksal/yasal durumlarla çelişmediği sürece istediğiniz yere gidebilirsiniz.
Ama fikirler de kimliğinizin en az bedeniniz kadar bir parçası, ve bunlar
üzerinde kontrol talep ettiğinizde ne yazık ki cloud sistemlerin size verdiği
yanıt “profilinizi görünümden kaldırdım, fakat kopyaları yedek olarak bende
durmaya ve ben bunları kendi sistemimde kullanmaya devam edeceğim” oluyor. Ama
tekrar ediyorum, bu veriler tamamen size ait ve yasal olarak bunlardan siz
sorumlusunuz.

Çözümüm var ama 2011’de de yazsam, 2014’de de yazsam umutlu değilim. Veri
lisanslama yöntemi gelmeli, ilgili kaynağım olsa bununla ilgili gerçekten
çalışmak isterim fakat akademik ve profesyonel uğraşları çıkarttığımızda kalan
kişisel boş zamanımı “özgür dünya”nın başka bir alanında, açık kaynak projeler
üretmekle meşgul ediyorum.

Daha kısa vadeli çözüm ise: yasaların bu verilerin de sahipliğini kesin olarak
bireylerin kendisine vermesi, bireyin kontrolü dışında işler yapan sistemlere de
dava açılabilmesini sağlamak olurdu. Yasa yapanlara ulaşmak ve bunu gündeme
aldırmak gerekiyor ama bu da bana pek gerçekleşmesi yakın bir olasılık gibi
gelmiyor. Şu web şirketlerini korumakta olan “ama anlaşmayı okuyarak siteye
girdin” çok zayıf bir bahane olmasına rağmen bir şeyler değişmedikçe ne yazık ki
geçerliliğini koruyor.
', 'tr          ', '01J6ZJ7S7MEQH7T95RG9NBBYFY', '', 'İçerik sahipliğine dair yasal eksiklikler');
insert into "story_tx" ("content", "locale_code", "story_id", "summary", "title") values (E'Uzun bir aradan sonra tekrardan blog yazmaya başladım. Bu esnada geçen süreye
KoçSistem’den ayrılarak Setur’a geçişimi, birçok yeni açık kaynak projeyi,
[7 farklı organizasyonda konuşma](../../../presentations/tr/README.md)yı ve bir
de yüksek lisans tezini sığdırabilmişim.

Genel olarak daha soft yetenekler, yazılım ürünleri, mimari, takım dinamikleri
ve pratikler üzerine yazmayı seviyorum. Bu nedenle geliştirici perspektifinden
çok daha farklı bir blog içeriğim olduğunu düşünüyorum. Ancak geçenlerde twitter
üzerinde geliştirici perspektifiyle yapmış olduğum bir React eleştirisini net
olarak ifade edebileceğim en uygun mecra da burası.

Eleştirim basitçe aşağıdaki tweet’de yer alıyor. Gelen yanıtlarla bir takım
tartışmalar oldu. Ben de sonunda çamur atıp kaçmış olmamak üzere bu yazıyı
planladım. Bugün ise ancak yazabiliyorum.

%[https://twitter.com/eserozvataf/status/1055860281425960961]

Bandı biraz geri saralım ve React’ın neden benim için fenomen bir view çözümü
olduğunu açıklamaya çalışayım.

jQuery gibi çözümler hem plugin host, hem polyfill, hem animasyon kütüphanesi,
hem DOM seçicisi ve manipulatörü, hem de bir web sayfasındaki akışı kontrol
etmeye çalışıyordu. Bu yaklaşıma alternatif olarak React, kapsam olarak yalnızca
kendisine verilen doküman yapısını (örneğin JSX) render etmekle sınırlı tuttu.
Bunu yaparken de bilgisayar biliminin neredeyse 1960’ların sonunda konuştuğu
bileşen-temelli UI yaklaşımından destek alıyordu. UNIX felsefesinin bir araçtan
beklentisinin özeti gibiydi, “bir iş yap, onu doğru yap”.

Çevremde daha fazla hayranı olsa da Vue için aynı şeyi söyleyemeyeceğim örneğin.
Çünkü Vue popülerite kaygısını belli edebiliyor. İlk sürümü angularjs yolundan
giderken, ikinci sürümü React, bir sonraki sürümü ise React Hooks’a doğru
evrilip amacını geliştiricinin konforunu arttırmak olarak belirleyebiliyor. Bu
bir tercihti ve React konfor yerine mühendislik doğrularını seçiyordu. En
azından ben öyle biliyordum.

React’ın yetersiz kaldığı yerde React-Router, Redux gibi diğer kütüphaneler
devreye giriyordu. Vue’nun isviçre çakısı olmaya çalıştığı bir gerçeklikte,
tornavida, çekiç ve çakı’nın montajda uygun yerlerde kullanılmasına
benzetebilirdik bu durumu.

Ancak, kullandığım ürünleri changelog seviyesinde takip eden biri olarak, React
16’dan itibaren projenin yolunu değiştirdiğine tanık olduğumu düşünüyorum.

Fiber güncellemesi (React’ın bir süre önce farklı bir rendering yaklaşıma
geçmesi) ile yalnızca Lifecycle’lar değişse, bu yeni rendering için bir takım
bağlantılar dışarıda kalsa bundan daha doğal bir şey olamazdı.

Bunun yerine, yeni tanıtılan her feature ürünün kapsamına (product owner’lar
neden bahsettiğimi daha net anlayacaktır) yapılmış bir eklenti olmaya başladı.

- getDerivedStateFromProps ile react-redux,
- yeni context api ile redux,
- Suspense ile [react-loadable](https://www.npmjs.com/package/react-loadable),
- useMemo ile memoize-\\* kütüphaneleri,

adeta ürünün temeline yerleştirdi ve/veya yerleştirilmeye başlandı. Tıpkı
vue’nun, vuex’e yaptığı gibi.

Şimdi bunların er ya da geç bir SPA’nın içinde mutlaka bulunması gereken
konseptler/bileşenler olduğunu, React’ın çekirdeğine yakın çözümlerin çok daha
güzel practiseler sunarak bunları geliştiricilere verebileceğini
savunabilirsiniz. Ama bu yukarıda çizdiğim React portesinin dışında bir React’ı
tarifliyor bana kalırsa. Kütüphane kimliğinden çatı/framework kimliğine doğru
kayacak veya opinionated olacak bir React hoşuma gitmezdi.

Twitter’da sonra geriye dönüp acımasız bir üslupla yazdığımı hissettiğim
“Geliştirilecek madde kalmamış, bunun yerine istekleri gerçekleştirmeye
başladılar” eleştirimin dayanak noktası da biraz buydu.

Yazının “plot twist”i ise burada. Ben React Hooks’un React değil de tamamen
farklı bir ürün olmasını savunsam da React Hooks’la yaptığım geliştirmeler ve
ufak projeler sonucu yeni React’dan çok memnun kaldım.

React, Pure Component yazabileceğimiz, fonksiyon ile tanımlanan componentlerin
100% stateless functional component olduğu, lifecycle eventler ile ilerleyen bir
yapı idi.

React Hooks ise functional componentlerin olduğu, “use” ile başlayan hook’lar
ile kullanıcıya side effect’li (dolayısıyla aslında stateless olmayan)
bileşenler kodları ürettiren apayrı bir yapı.

Fakat günün sonunda React Hooks ile yazılan yapıda,

- Prop’ları paslamanın, state yönetmenin daha kolay olduğu,
- Redux’a ve React-Redux’a ihtiyaç bırakmayan,
- Fonksiyon olarak bileşen tanımlamanın lifecycle’larla yazılmış bir sınıftan
  kolay olduğu,
- Lifecycle eventler yerine hook’ların kullanımı ile temiz, anlaşılır ve
  okunabilir kodların oluştuğu,

bir kod tabanı vaad ediyor.

React Hooks ile yazılmış bir örnek için,
[kişisel web sitemin ana sayfasının kodu](https://github.com/eser/eser.ozvataf.com/blob/archived/old-site/src/app/pages/frontpage/index.tsx)nu
inceleyebilirsiniz.

Son olarak birkaç saat önce Armağan Amcalar ve Fatih Kadir Akın ile Kod
Analizi’nin ilk konuğu oldum. YouTube üzerinde canlı yayında
[Evangelist](https://github.com/eser/evangelist) isimli açık kaynak projemi
ellerinde küreklerle gömdüler. Daha sonra farklı bir proje ve soru/yanıt ile
devam ettik. Hem tarihe not olması, hem de Kod Analizi’nin kitlelere ulaşması
için buraya eklemek istedim.

%[https://www.youtube.com/watch?v=e-1aIc93QXg&list=PLWLiJPAYmgZAtN-trvpEpKysVL2LKMcO6&index=1]

React ve React Hooks üzerinden tartışmalar bana kalırsa henüz yeni başlıyor. Ben
de birçok fikrimi elbette değiştirebilirim. Bu yazıda umarım biraz da olsa
meramımı anlatabilmişimdir.

Sağlıcakla kalın.
', 'tr          ', '01J6ZJ7SXTDJA9SZ9WGV7RD03H', 'React''ın kütüphane kimliğinden framework''e kayışına dair eleştirilerimi anlattım, ancak React Hooks''un sağladığı pratiklik ve kod temizliğinden oldukça memnunum.', 'React’ın Vizyonu ve Hooks');
insert into "story_tx" ("content", "locale_code", "story_id", "summary", "title") values (E'Yazının ilk bölümü olan
[Uzaktan Çalışma Deneyimim  —  İşverenler ve İş Arkadaşları (1/2)](./20161010-uzaktan-calisma-01)nda
tam bir “kişisel weblog” hikayesi sunarak kendi uzaktan çalışma deneyimim ne
zaman gerçekleşti, neler oldu ve nasıl sonuçlandı üzerinden konuyu aktarmaya
çalıştım.

Önceki yazının okuyucu açısından pek ortak nokta bulunmadığı sürece okuması
sıkıcı bir yazı olduğuna katılıyorum. Neyse ki bu yazıda ne zaman, neler ve
nasıl oldu sorularının dikte ettiği hikaye kısmından sıyrılarak direkt olarak
deneyimin kendisine odaklanacağım.

#### Avantajlar

Öncelikle uzaktan çalışmak kişinin kendisi için mükemmel bir deneyimdir. Eğer
denge sağlayabilirseniz uzun vadede:

- Ancak kendinizi verimli hissettiğiniz saatlerde çalışan ve bu sayede büyük
  orantılarda verimli çalışan,
- Hobilerine, ailesine, çevresine, sağlıklı ve mutlu olmaya zaman ayırabilen,
- Vakit kaybetme alışkanlığı olmayan,
- Yol, sokak, hava durumu stresi çekmeyen,
- İstediği yemeği/içeceği istediği lezzette ve sağlıklı bir şekilde temin
  edebilen,
- Gürültü, kalabalık, konsantrasyon dağıtan etkenlere uzak olan,
- Hem işi hem özel hayatı dengeleyebildiği için iç huzuru ve sorumluluk
  duygusunu geliştirmiş olan,

bir insan haline geliyorsunuz.

#### Engeller

Fakat her şey göründüğü gibi toz pembe değil, işyeriniz “uzaktan çalışma
imkanını” sağlamaya ikna olduğunda dahi bunun sizin için “sonun başlangıcı”
olmaması için önünüzde bazı sınavlar var.

Öncelikle bunları nasıl aşabileceğimize bakmak için “reverse engineering”
yaparak olayı tersten ele almamız, işyerinden çalışmanın avantajlarını SWOT’a
tabi tutup bunu “uzaktan çalışma şartlarımız”la dengelememiz gerekmekte. İşte
başlıyoruz:

**İşyerinde çalışma**

- Ofiste olmanız gereken saat genellikle belirlidir.
- Etrafınızdaki herkes “mesaide”dir ve çalışmaktadır.
- Tüm çalışma ortamı tek bir amaç için tasarlanmıştır “çalışmak”.
- Çevrenizde farklı bir iş yapmaya başladığınızda, uzun süre ortadan
  kaybolduğunuzda sosyal olarak sizi ayıplayacak hatta sizi cezalandırabilecek
  iş arkadaşlarınız vardır.
- Çalışma arkadaşlarınızla yüzyüze olduğunuz için iletişimler daha empati dolu,
  daha sıcak ve daha samimi gerçekleşir.

**Uzaktan çalışma**

- Çalışma ve iş-dışı zamanınızı kolaylıkla karıştırabilir, çalışmanız gereken
  saatlerden veya kendi zamanınızdan ödün vermeye başlayabilirsiniz. Örneğin
  aynı anda hem çalışıp hem başka bir şeyler yapmaya çalışırsanız (ev işi,
  televizyon izlemek, ortamda bulunan veya internetten arkadaşlarla sohbet
  etmek) muhtemelen ikisini de başaramayacak, toplamda bu durumdan zararlı
  çıkacaksınız.
- Tek başına olduğunuzdan “birden” çalışmaya başlamak oldukça zordur, gerektiği
  zamanda bilgisayar başında olsanız dahi “şu diziyi de izleyip öyle çalışmaya
  başlayayım” diyerek procrastination’a teslim olabilirsiniz.
- Evdeyseniz kendi yatağınız, rahat koltuğunuz, buzdolabınız, kahve makineniz,
  televizyonunuz, müzik setiniz, netflix aboneliğiniz sizin kendi kendinize
  kurduğunuz tuzaklardır. Verdikleri konfor ve rehavet sizin “düşük tempoda”
  çalışmanız, sürekli bölünmeniz anlamına gelebilir.
- Çalışma arkadaşlarınız ve müşterileriniz ile yazılı iletişim üzerinden
  haberleşirsiniz ve bu özellikle işler kötüye gitmeye başladığında iki tarafın
  da birbirini gözlemleyemediği, göz teması kuramadığı bir ortamda alabildiğine
  sert üsluplar içeren e-mailler, uyarılar ve hatta suçlamalara dönüşmeye
  oldukça meğillidir.
- Düşük olasılık gibi görünse de “işinizden” dolayı para kazandığınızı siz veya
  çevreniz unutabilir. İnsanlar sizi çalışmıyormuşsunuz gibi değerlendirebilir,
  işinize ayırmanız gereken zamanı koruyamayabilirsiniz.

#### Çözümler

İş ortamının bize sağladığı avantajları ve uzaktan çalışmaya dair olası
tehditlerimizi analiz ettiğimize göre şimdi aradaki farkı nasıl
kapatabileceğimiz üzerine kafa yormaya başlayabiliriz.

Öncelikle yapmanız gerekenlerin bir listesini oluşturalım:

- Kişisel motivasyon ve disiplin sağlamalısınız.
- Çalışırken odaklanmalı, zamanı boşa harcamamalı, iş takibi yapabilmelisiniz.
- İş ve özel zamanınızı korumalısınız.
- Sosyal olarak çevrenizi çalıştığınıza ve meşgul olduğunuza ikna etmelisiniz.
- İş arkadaşlarınız ve müşterilerinizle iletişimi canlı ve pozitif tutmalısınız.
- İş yapmaktan kazancı olan bir profesyonel olduğunuzu asla unutmamalısınız.

Peki bunları nasıl gerçekleştireceğiz, işte başlıklar halinde benim edindiğim
bazı prensipler ve yardımcılar:

**Organize Olmak**

Mesaimizi bitirirken bir sonraki çalışmamızın ne zaman başlayacağını ve o
“oturum”da neler yapacağımızı ve ne kadar çalışacağımızı önceden planlamamız
gerekiyor.

Organize olduğumuzdan, bir sonraki çalışmaya başlayacağımız zamanı da biliyoruz.
Bu sayede iyi bir çalışma geçirmek için ön hazırlıklar yapabilecek konumdayız.
Bu da bizi bir sonraki maddeye hazırlıyor.

**İyi Dinlenmek ve Çalışma Öncesi Zinde Olmak**

Çalışırken iyi hissetmek önemli bir etken. Bitkin bir halde çalışmak, o mesai
içerisinde istediklerimizin hiçbirini gerçekleştiremeden çalışmayı sonlandırmak
olmak maalesef ki o gün çalışmamış olmakla eşdeğer. Bu nedenle,

- Gerekirse erken yatmak ve sağlıklı bir uyku uyumak,
- Alkol veya sakinleştirici ilaçlar gibi düşünce hızını veya metabolizmamızı
  değişikliğe uğratan etkenlerin kullanımında düzenlemeye gitmek,
- Bol bol oksijen almak,
- Gürültüden uzak durmak,
- Spor veya yürüyüş yapmak,

zinde hissetmenin önemli anahtarları oluyor.

**Mekan değiştirmek**

“Evden çalışmak”, uzaktan çalışmanın ikinci safhası olarak değerlendirilebilir.
İlk etapta kendinizi iyi disiplin edemediyseniz “evden çalışmak” değil “uzaktan
çalışmak” üzerine yoğunlaşabilirsiniz.

Günümüzde coworking space’ler, Starbucks şubeleri, bazı AVM’lerin avluları hem
rahat bir ortam, hem internet hem de elektrik sağlayarak sizin mesai saatlerinde
zaman geçirebileceğiniz bir ortam sağlayabiliyorlar. Ev yerine bu tarz
mekanlarda olmak, çevrenizde insanların kendi halinde olduğu, sizin de orada
bulunmanızın bir amacı olduğu durumunu oluşturduğu için motive edici olabiliyor.

Bu motivasyonu evde de imkanınız varsa ayrı bir çalışma odası oluşturarak ve
çalışma saatlerinde oturma odasına uğramayarak oluşturabilirsiniz. Tabii ki yine
kendi disiplininizi henüz sağlayamadıysanız o oda içerisinde film izlemeyi,
telefon görüşmesi yapmayı, dikkat dağıtıcı şeyleri kendinize yakın tutmayı
yasaklamak iyi bir başlangıç olacaktır.

**Zamanı takip etmek**

İster
[(10+2)\\*5 Procrastination Hack](http://www.43folders.com/2005/10/11/procrastination-hack-1025)
veya [Pomodoro Technique](http://pomodorotechnique.com/) gibi teknikler, ister
[RescueTime](https://www.rescuetime.com/) gibi araçlar kullanın amacınız iş ve
özel zamanı birbirinden ayırmak olmalı. Araçlar en başta sizin dikkat
dağınıklığınızı engelleyip verimli çalışmanızı, gerektiği kadar sizi dinlenip
zaman kavramını daha da özümsemeniz konusunda yararlı oluyor.

Böylece zaman içerisinde işlere daha net zaman koyabilir, ne kadar çalışmanız
gerektiğini daha net kestirebiliyor oluyorsunuz.

**Mümkün olduğunca fazla toplantı ve iletişim**

İş arkadaşları ve müşteriler her zaman işimizin bir parçası ve biz temelde
ürettiğimiz “yazılım çözümleri”nde onlarla birlikte aynı cephede “bilgisayara
isteğimizi anlatmaya” çalışıyoruz. Bu noktada bizim bireysel olarak bilgisayarla
iletişimimizin son durumunu cephe arkadaşlarımızla paylaşabiliyor olmak,
gerçekleştirdiğimiz işin hesap verilebilirliği açısından büyük önem arz ediyor.

Bu noktada birçok kişi “yahu ne gerek var paso toplantıya, skype’dan görüşürüz
işte gerekirse” görüşünü savunacak olsa da Agile Manifestosunun da önerdiği gibi
işimizin odağında insan olduğunu unutmamak ve sıcak iletişimin bizi hem üzerinde
uğraştığımız projeye hem de çalışmaya yönelik motive edeceğini bilmek gerekiyor.
Amacını aşmayan, iletişimi güçlendiren her toplantı ve bir araya gelme
organizasyonu sizi ve işin içine dahil olan herkesi projeye sıkı sıkı
bağlayacaktır. Yazılı iletişime önem veriyorsanız dahi bir kahve partisi verin,
geliştirirken yaşadığınız zorlukları paylaşın yeterli olur.

**Rapor Yazmak**

Kimseyle paylaşmayacak dahi olsanız o hafta hangi işler için çalışmanın başına
oturduğunuz ve neleri tamamlayıp çalışma ortamından ayrıldığınızın bir raporunu
oluşturun. Bunun için Trello tarzında Kanban boardlar da kullanabilirsiniz,
stickerlar da not defteri de…

Yine de iş arkadaşlarınız istemese dahi bu raporları onlarla paylaşmanın
iletişimi güçlendireceğini de eklemek gerekiyor.

Benim deneyimlerim genel itibariyle bu şekilde. Yazıda mümkün olduğunca farklı
konuları aynı başlıkların altında toplamaya çalıştım. Fakat buna itirazlar
gelebilir veya ben fikrimi değiştirebilirim, bu nedenle yazıyı geliştirmek adına
daha farklı zamanlarda güncellemelerde bulunabilirim.

Firmalar gün geçtikçe daha fazla uzaktan çalışma modeline sıcak bakmaya
başlıyorlar. Ben de deneyimlerimi aktararak umarım bu dönüşümde ufak da olsa
soruları yanıtlayarak bir katkıda bulunabilmişimdir.
', 'tr          ', '01J6ZJ7SGH8HBPKXMWD3G0JGW9', 'Önceki yazıda ele aldığım uzaktan çalışmanın kişisel deneyimlerinden sonra; bu yazıda avantajlarını, zorluklarını ve çözümlerini ele aldım. Disiplin, zaman yönetimi ve etkin iletişimin kritik olduğunu, çalışma ortamını doğru seçmenin ise motivasyonunu aktarmaya çalıştım.', 'Uzaktan Çalışma Deneyimim  —  Kişisel Organizasyon (2/2)');
insert into "story_tx" ("content", "locale_code", "story_id", "summary", "title") values ('[Kariyerimin 21. yılı](../../../curriculum-vitae/en/README.md)ndayım. Kendimi
Türkiye sınırları içerisinde “büyürken evinde bilgisayar olan, bu sayede yazılım
teknolojisinin gelişmesine tanıklık etmiş ilk neslin az sayıda üyesinden biri”
olarak tanımlıyorum (bkz:
[90’larda programlama öğrenmek](../20160221-90larda-programlama-ogrenmek)).
Dolayısıyla şartlar ilgim ve isteğimle birleşince benim için doğal bir sonuç
oldu. İkinci tercihim Bilim Kurgu Yazarı olmaktı, ama o zaman işim daha zor
olurdu.

## Geçmiş

#### Meslekteki İlk Yıllarım

Yazılım geliştirici olarak kariyerimin başlangıcında, bugünkü “birden fazla
insan tek ürün/projede çalışır” anlayışının tam tersine işleyen “bir insan
birden fazla projenin tek çalışanıdır” düzeninde çalışıyorduk. Mesai
arkadaşlarımın “presales” faaliyeti olarak müşterilere gitme, toplantılara girme
faaliyetlerini hiç sevmiyordum hatta aynı ofisi paylaştığım insanlara “günaydın”
demek konusunda bile isteksiz olan, dışarıdan bakıldığında bariz bir şekilde
iletişim konularında sınıfta kalan biriydim.

Zamanla aldığım sorumlulukları yerine getirmek, başarılı olmak, zorlukları
yenmek motivasyonum beni değiştirdi. Artık yazılımla sunabileceğim “doğru iş
çıktıları” kendi doğal sınırlarını bulduğunda insanlarla çözümler arıyordum.
Doğru iletişim dili, müzakere etmek gibi yetenekler kazanarak bugünkü kişiliğime
evrimleşmiş olsam da, ben o günlerdeki yaklaşımlarımı ve davranışlarımı hiç
unutmadım.

Sonrasında 12 yıllık bir süre zarfında bazen startup’larda, bazen yazılım
evlerinde backend yazılımcı, senior full-stack yazılımcı, ekip yöneticisi, head
of development gibi pozisyonlarla çalıştığım bir süreç oldu. Yüksek lisans için
tekrar kolları sıvayıp bir yandan uzaktan danışmanlık verdiğim, bir yandan da
araştırma görevlisi olduğum bir sürece girene dek bu tempo ile devam ettim.

#### Koç Holding Çatısı Altında Çalışmak

Yüksek lisansı bitirdikten sonra İzmir’de hayat kurmaya çalışırken, bir yandan
İstanbul’daki müşterilerim için haftanın 4 günü seyahat halinde olan, hands-on
danışmanlık veren bir hayat düzenindeydim. Yerleşik hayata geçiş hayaliyle
“kurumsal hayat” maceram başladı.

Koç Holding bünyesinde KoçSistem’de Arçelik, Tüpraş, Koç Üniversitesi v.b. iyi
profilli firmaların projelerinde hizmet yöneticisi olarak çalıştım. Daha sonra
yeni bir yapılanma içerisinde olan Setur’a transfer oldum. Yerleşik hayata geçme
planları yaparken, aksine seyahatlerim arttı. Artık şirket genel merkezi
yanısıra müşterilere de seyahatlerim vardı.

Bu dönem aslında teknoloji sektörü de bir dönüşüm içerisindeydi. Startup’ların
ivmesine yaklaşmaya çalışan kurumsallar, internetin hızlanması, bilginin daha
kolay yayılması, mobilitenin artması gibi faktörlerle birlikte bulut bilişim,
yönetilen hizmetler, dağıtık sistemler, yeni platformlara geçmeye istekliydiler.
Ama önlerinde birçok engel vardı.

En basitinden:

- gantt chartları yok edip değişen taleplere çevik tepkiler verebilmek,
- ürün ekipleri oluşturarak knowhow’ı tutmak, optimal çözümleri tartışıp
  uygulamak,
- .NET gibi monolitik yapılardan .NET Core gibi modüler yapılara geçiş,
- CI/CD ile çalışmak, kesintisiz deployment’lar yapmak,
- statik kod analizleri yapmak, güvenlik ve potansiyel zaafiyet raporları almak,
- kütüphaneleri belirli package repository’lerde tutmak,
- entegrasyon ve davranış testleri yazmak,
- uptime’ı yüksek tutacak şekilde servis orkestrasyonu yapmak,
- tasarım dili oluşturup UX tasarlamak,
- SVN/TFS gibi yapılar yerine Git gibi distributed source control yapıları
  kullanmak,

v.b. hedefler vardı. Ama buralara nasıl varılabilirdi?

Kurumsaldaki hayatım öncesinde danışmanlık verirken benim alışkın olduğum bu
yapılar, kurumsaldaki dünyanın pratik etmediği disiplinler içeriyordu. Aynı
zamanda kişilerin alışkanlıkları yıkarak sıfır noktalarına inmeye cesaretleri
veya antremanları olmadığı için kurumsal dünyada yıkılamaz birçok sabit vardı.

Pandemi’ye kadar süren uzun bir maraton sonrasında, her zaman önümü açan
yöneticilerimin de desteğiyle, son geldiğimiz nokta flat bir organizasyonda bana
bağlı 9 agile ürün ekibinin multi-disipliner (ux, qa, fe, be, analiz), polyglot
(birden fazla programlama dili ile) ve otonom çıktılar verdiği, kubernetes’de
kendi servislerini yönettikleri bir noktaya gelmiştik.

Bu süreç içerisinde zamanımı %20 bilgisayar başında, %20 ekiple birlikte olmak
için İzmir-İstanbul arası yolda, %60 ise insanlarla toplantı odalarında
geçirdim. Edindiğim know-how’lar çok işime yaradı.

İnsanları bürokratik, eskimiş, angarya-dolu bir yapıdan kurtarmış olup mesleki
anlamda yazılım mühendisliğinin modern doğrularına yaklaştırmış olmanın manevi
tatminini bolca yaşadığım bir dönemdi. O dönem uyguladığım iletişimde de
başarılı olduğumu, kötü bir yönetici olmadığımı, halen o dönem çalıştığım
arkadaşlarımla halen temasta/paylaşımda olmamdan anlayabiliyorum.

## DevRel Alanında Çalışmak

#### DevRel Danışmanlığı

Yukarıda seyahatsiz, yerleşik bir düzen arzuluyordum yazmıştım. Ama ne
istediğime dikkat etmeliymişim. Pandemi’den sonra ekiple birlikte olmak için
İzmir-İstanbul arası seyahat şansım kalmamıştı. Hatta Setur gibi ana işi turizm
olan bir firmanın hem rezervasyonlarını iade etmesi, hem de pandemi koşullarıyla
tekrar ne zaman hareketleneceğinin belli olmaması durumu bizi derinden sarstı.
Belirsizlik nedeniyle yatırım yapamaz hale geldik ve elimizdeki yetenekli
kadroyu yavaş yavaş sektörde iyi teklifle gelen diğer firmalara kaybetmeye
başladık.

Çok iyi scale-up ettiğimiz bir dönemin ardından, istemsiz lay-off yaşadığımız
bir döneme girmiştik. Bir önlem alamayacağımızla barıştığımız noktada da ben de
CIO’muzla konuşarak, dövizle gelen yurt dışı tekliflerini değerlendirip
danışmanlık hayatıma geri dönmeye karar verdim.

Danışmanlığa geri döndüm, şirketleştim. Hali hazırdaki hands-on development
işlerimi paylaşmak için beraber çalıştığım arkadaşlar olmaya başladı. Vergi
mükellefliği ve işverenlik derken kurmuş olduğum “Açık Yazılım Ltd. Şti.”nin
sektörde olmayan, değerli bir alanda müşterilerine nasıl bir “danışmanlık ürünü”
sunabileceğine kafa yormaya başladım.

Hands-on programlama hizmeti vermek “fark yaratan” bir hizmet olmayacaktı.
Ekiplere Software Craftsmanship’lerini ve çevikliklerini arttıracak şekilde
koçluk sağlayabiliyordum, ürün ekiplerinin kurulması ve DevOps pratiklerinin
benimsenmesi konusunda da... Ancak bunlar halihazırda organizasyonel koçluk
veren birçok firma tarafından sağlanıyordu, dolayısıyla yine başkalarının hali
hazırda verdiği bir hizmetten bahsediyor olacaktık.

Sherlock hikayelerinde geçen ünlü bir söz vardır: “İmkansızın gerçekleşme
ihtimalini elediğinizde, olasılıklardan geriye kalan ne kadar imkansız olursa
olsun, gerçeğin kendisi olmalıdır”.

Ben de benzer bir eleme yapıp yöneticilik günlerimde sırtlandığım konular
arasında hali hazırda sektörde fark yaratacak kadar eşsiz verebileceğim
hizmetleri sıraladım: Developer Experience, Employee Experience, Talent
Acquisition, Community Management, Engineering Success, Tech Branding, Open
Source gibi konular ortaya çıktı.

Elbette ki yüksek lisans zehirlenmesi olan “literatür taraması”nı gerçekleşip
daha önce de birkaç kez duyduğum ama derinlemesine araştırmadığım DevRel, yani
“Developer Relations” kavramına denk geldim. Okudukça konu ile ilgili
bilgileniyor, ilk düşündüğüm gibi Developer Advocacy kadar sınırlı olmayan hatta
Tech Marketing’den de geniş bir alan olduğunu fark ediyordum.

Potansiyel müşterim olabilecek şirketlere DevRel hakkında sunumlar hazırlamaya
ve sunmaya başladım. Artık verdiğim danışmanlığı DevRel danışmanlığı olarak
isimlendiriyordum ve hatta bir gece pek değerli Furkan Kılıç ile yaptığımız
geliştirici odaklı spaces yayınlarından birini kapatırken “DevRel diye bir
disiplin var, ülkemize daha uğramadı ancak bir gerçek, ileride bir programda
bunu da konuşalım” demiştim.

DevRel danışmanlığı için 3 firma ile bu esasta görüşmeler halindeydim ki; hali
hazırda danışmanlık verdiğim Getir, agresif bir işe alım hedefi olduğunu,
teknoloji kadrosunu global’de 1000 kişiye genişletmeyi planladıklarını ileterek
benden bir çözüm beklediler. Elbette bu Getir büyüklüğünde “Getir’e dışarıdan
destek verilerek” yapılabilecek bir sorumluluk değildi; ben de bu disiplini ülke
sınırlarında ilk kez uygulayacak olma motivasyonuyla hiçbir detayla ilgilenmeden
Getir’in “DevRel Director”u olmayı kabul ettim.

Yazının devamı için:
[DevRel (Developer Relations) Yolculuğum — DevRel Director Olmak](./20221005-devrel-02)
', 'tr          ', '01J6ZJ7N7VSV25K6ZD7ZVMQ0KV', 'Yazılım kariyerimin ilk yıllarında iletişim becerilerim zayıftı ama sorumluluklarım sayesinde bu yetkinlikleri geliştirdim. Zamanla farklı rollerde çalışarak, yazılım mühendisliğinin modern pratiklerini kurumsal dünyaya taşıdım. Bu yazıdaysa pandemi sonrası DevRel kavramıyla tanışmamı ve Türkiye’de bu alanda öncü rol üstlenmemi anlatıyorum.', 'DevRel (Developer Relations) Yolculuğum — Nasıl başladım? (1/2)');
insert into "story_tx" ("content", "locale_code", "story_id", "summary", "title") values (E'6 Şubat 2023’de yaşadığımız coğrafyada elim bir olay gerçekleşti.
Kahramanmaraş’ta gerçekleşen ve en az 10 ilimizi daha etkileyen deprem felaketi
sonucunda aynı coğrafyayı paylaştığımız insanlar kayıplar ve travmalar yaşadık.

Daha önce herhangi bir arama-kurtarma eğitimi almamış, bölgeye gitmek istese
yolda durdurulacak, örgütsüz ve bir şeyler yapmak için gönüllü olabilecek
insanların ilk aklına gelen şeyi gerçekleştirmek istedik: profesyoneli olduğumuz
konular üzerinden felaketten etkilenenlere yardımcı olmak.

Bu acı olay vesilesiyle, son kontrol ettiğimde nüfusu 23,950 kişiyi aşmış,
yazılım alanında faaliyet gösteren gönüllüler olarak bir Discord sunucusunda bir
araya geldik ve kendimizi Açık Yazılım Ağı olarak isimlendirmeye başladık.

**Açıklama:** 6 Şubat’tan bu yana zamanı efektif kullanmak için hep kısa kısa
yazdım, kimsenin zamanını fazladan tüketmemeye çalıştım. Discord’da dahi sahnede
en az benim konuştuğumu, twitter’da oldukça az tweet attığımı
gözlemlemişsinizdir. Bugün bu uzun yazıyı kaleme almamın nedeni insanların
bundan sonrasına ilişkin özelden ilettikleri mesajları teker teker değil, toplu
halde yanıtlamak. Ve halen iletişim yapamadığımız konular, soru işaretleri varsa
açık bir biçimde yanıtlamak için bu yazıyı yayınlıyorum.

## Nasıl toplandık?

Hem daha önce topluluk adına birçok yayın yaptığım hem de mesai arkadaşım Furkan
Kılıç’ın Twitter’da yine yazılım ekosisteminden tanıdığım İrem Kömürcü ile
mesajlaşmalarını gördüm, nasıl yardım edebileceklerini tartışıyorlardı. Ben
şahsen henüz hiçbir haber takip edememiş, şirketin planlama toplantımıza
hazırlanıyordum. Toplantımız esnasında bir yandan Twitter’daki mesajları göz
ucuyla okuyor, Furkan’ın ortaya attığı “bir ısı haritası oluşturalım” fikri için
açılan Telegram grubunu takip etmeye çalışıyordum.

Toplantımız biter bitmez Discord’a geçtik. Furkan’la “benim genel yetkilendirme
ile ilgileneceğim, onların uygulamaya başlayacağı” ufak bir planlama yaptık.

Furkan’ın “sosyal medya’da insanların ürettiği yardım çağrılarının ısı
haritasını oluşturma fikri” o kadar temel bir doğruya dayanıyordu ki... Makro
ölçekte ülkede yetkililer dışında kalan hiç kimsenin bölgeden haber alamaması,
mikro ölçekte ise olayın canlılığı ile daha mahallesinde ne olduğunu
kavrayamamış insanların hangi bölgede nasıl bir yoğunluk olduğunu
gözlemleyebilmesi sağlanıyordu. Bu basit ama güçlü fikri “gelin birlikte
geliştirelim” çağrısı ile dalga dalga yayıldı. İnsanlar “IT Deprem Yardım”
ismiyle oluşturduğumuz Discord sunucusunda toparlanmaya başladılar.

## Nasıl büyüdük?

Bilişim sektörüne temas eden herkese yaptığımız açık çağrı yanıt bulunca
Furkan’la ilk stratejimiz, onun projeyle ilgilenmesi, benimse biraz daha genel
resme bakıp yetkilendirmeleri yapmak, repository’leri açmam üzerine oldu.

Kendi adıma yapacaklarım zaten kafamda netti,
[github.com/acikkaynak](https://github.com/acikkaynak) 2015 yılından bu yana
vardı ve zaten tam da insanlarla birlikte geliştirme yapmaya hazırdı. Kendi
çevremde, topluluktan tanıdığım arkadaşlarımızla da yürütülen ufak projelerde
benimsenen fikirsel altyapı birçok profesyonelin de bu çalışmalara dahil
olabileceği bir zemine elverişliydi.

Furkan ise frontend, backend, tasarım gibi ekipler oluşturup bu kişilerden
bazılarını projelerde yetkilendirilmek için bana yönlendiriyordu.

Projelerimizi o noktadan itibaren Apache 2.0 açık kaynak lisansı ile herkesçe
denetlenebilir, herkesin nezaretinde geliştirilen bir yapıya kavuşturduk. Veri
dışında her konuda şeffaf olduk, verileri yalnızca STKlar/yetkili
kurumlar/akredite kurumlar ile paylaşacak şekilde bir ilke belirledik. Bu bizim
süreç boyunca hassasiyetle koruduğumuz temel iki ilkemizi oluşturdu.

Son kontrol ettiğimde yalnızca discord sunucumuzda 24,000 kişi olmamıza az
kalmıştı. Süreçte mutlaka küstürdüğümüz, kendini dahil olamamış hissederek
uzaklaşan kişiler de olmuştur.

![](https://cdn.hashnode.com/res/hashnode/image/upload/v1676612473943/89951b44-d922-4e81-a0cb-9404d16c9701.png
align="center")

## Açık kaynak olmasının anlamı nedir?

Projelerin açık kaynak olmalı kriterini benimsedik çünkü başından bu yana
amacımız: yaptığımız iş denetlenebilir, gerektiğinde belirli sistemlere entegre
edilebilir ve diğer insanların da az bir eforla bir ucundan tutabileceği şekilde
geliştirmeye açık olmalıydı.

Kurumsal bir şirkette bir projeye dahil olan birileri, belki 3. haftasının
sonunda “süreçlere uygun” geliştirme yapmaya başlayabilir durumda oluyorsa bizde
bu süre 5-6 dakikaya indi.

O an “projeyi çıkartma” motivasyonuyla düşünemediğimiz ama sonradan da fark
ettiğimiz bir nokta ise; gün 1’den açık kaynak ilerleyerek aslında ileride
bizden bağımsız bu sistemlerin herhangi bir ülkedeki başka herhangi bir
STKsının/yetkili kurumunun/akredite kurumunun bunları kullanabilmesini sağladık.

Açık kaynak olarak geliştirdiği için artık bizim bu projelerin kaynak kodlarını
“kapatma” veya başka bir yerde kullanılmasına engel olma şansımız yok. Bu
anlamda projenin kaynak kodları kimsenin sahipliğinde veya tekelinde değil. Kötü
bir vesileyle de olsa, bizim bu çalışmaları “gönüllüler” olarak yapıp, dünya
bilgi hazinesine bağışlamış duruma geldik.

## Ne tür projeler geliştiriliyor?

İlk olarak [afetharita.com](https://afetharita.com) ile çıktık. Daha sonra yeni
fikirler geldikçe yine açık kaynak bir yazılım vakfı gibi (Apache Foundation,
CloudNative Computing Foundation, Free Software Foundation, Linux Foundation,
Mozilla Foundation v.s.) gelen projeleri değerlendirip, bunlara gönüllü ekipler
kurup ilerlemelerinin önündeki engelleri kaldırmaya çalıştık.

Bir proje listesi üzerine çalışıyoruz, çok yakın zamanda iç süreçlerimiz daha da
düzene girmiş olacak. Ancak bugünlük “yayına hazır” projeleri
[https://afet.org/](https://afet.org/) adresi altından görüntüleyebilirsiniz.

## Denetim ve güvenlik ne aşamada?

Konu açık kaynaklı geliştirme olunca akla gelen ilk sorulardan birisini
eleyelim: Hayır, herkes dilediğince müdahale edemiyor. Öncelikle proje liderleri
projelerde yetkili oluyor. Ancak herkes geliştirilen kodları okuyup, yeni
değişiklik önerilerinde bulunabiliyorlar. Bunu bir meclis yapısı gibi
düşünebilirsiniz. Nasıl kanunlar herkese açıksa ve ve her kanun önerisi kabul
edilmiyorsa, bu süreç de “öneriler” ile gidiyor. Projede yetkili olanlar da bu
önerileri kabul edip, kendi şerhlerini düşüp projeye dahil etme veya etmeme
tasarrufunda bulunabiliyorlar.

Yine süreçte gönüllü olarak gece gündüz çalışan bir altyapı ekibimiz oluştu ve
onlar da bize süreç boyunca desteklerini esirgemeyen bulut sağlayıcılarla
birlikte güvenilir bir ortam sağlıyorlar.

Özellikle afetten sonraki 2. gün siber saldırıların da gelmesiyle yine sektörden
çok deneyimli siber güvenlik firmalarının ve uzmanlarının destekleri ile hem kod
tarafındaki hem de altyapı tarafındaki zaafiyetlerimizi belirli aralıklarla
kontrol ettirdiğimiz bir yapıya kavuştuk.

## Kimlerden destek aldık?

Neredeyse herkesten!

Gün 0’da gönüllülerden.

Saatler sonra bilinen bulut hizmetleri sağlayan tüm firmalar bize yardımcı
olmaya çalıştı. Hiçbir hizmet talebimizde hayır denilmediği gibi, mevcut
hizmetlerindeki limitleri bizler için kaldırdıkları oldu.

Proje ekiplerindeki insanların firmalarına haber vermesiyle birlikte firmalar da
bize ulaşıp “nasıl yardımcı olacaklarını” sordular. Birçok firma çalışanlarının
açık kaynaklı olarak katkıda bulunabilmelerini desteklediler. İlgili alanlarda
nasıl destek verebilecekleri konusunda kafa yorduk ve iletişime her zaman açık
haldelerdi.

## Sektörel kazanımlarımız neler oldu?

Bu afet durumu ile bir araya gelmenin ülkemize, coğrafyamıza çok fazla kayıp
yaşattığı aşikar. **İşin sahadaki STKlara/yetkili kurumlara/akredite kurumlara
geri hizmet sağlamak kısımlarını dışarıda bırakarak** (diğer türlü haddimizi
aşıp, diğer kurumlarla yarıştığımız anlamı çıkabilir, o nedenle lütfen maruz
görün) Tüm bu kayıplarımızı aklımızda ve kalbimizde tutarak, yalnızca buradaki
çalışma kapsamında **teknik** bir değerlendirme yapmamız gerekirse:

- Bilişim sektöründeki kişiler twitter’daki hararetli tartışmaları,
  küskünlükleri gerilerinde bırakarak, zor anlarda yardımlaşma ve dayanışma
  zemininde bir araya gelip, birlilkte çalışabildiklerini hissettiler

- Sahaya nasıl yardım ulaştırabileceği konusunda sabahlara kadar tartışan
  yüzlerce kişi hepimizin birbirine inancını tazeledi

- İnsanlar açık kaynak geliştirme, bulut hizmetleri v.b. yeni çalışma
  disiplinleri ile tanıştılar

- Birçok organizasyonun seneler içinde alacağı yolu sektörün belirli
  pratikleriyle 2-3 günde bir düzene getirebildik

- Kaos gibi görünen bir yapıda çalışan, canlıda olan yazılımlar örebildik

- Kurumlara, geçmişte birçok kez “hukuki” nedenlerle çekince ile yaklaştıkları
  “açık kaynak geliştirmenin” çok kaygılanacak bir yönü olmadığını ispat ettik

- Farklı kurumlar altındaki insanlar birbirleri ile çalıştılar, farklı
  yaklaşımlar, farklı kültürler bir araya geldiler ve birbirlerini tanıyıp yeni
  ilişkiler başlattılar

- Yazılım alanında sürekli tartışan insanlar gibi görünüyorduk, hem kendimiz hem
  de dış partiler potansiyelimizi gözlemleyebilmiş oldu

- Kurumlara, medyaya, çevremize, bizden sonraki nesillere ve en çok da kendimize
  ilham “influence” oluşturduk. Bir arkadaşımla yaptığım toplantıda 1-2 aydır
  yazılım alanında kötü zamanlar geçirdiğini ama buradaki birlikle tekrar
  doğduğunu konuştuk.

  Dahası şöyle mesajlar aldık: “Hatay’da 12 yaşında öksüz kalan akrabam bana
  _abi çocuklar için de bölüm ekleyelim bana bilgisayar alalım kod yazcam_ dedi,
  _[afetcocuk.com](http://afetcocuk.com)’u bana al_ diyor.”

## Bundan sonra ne olacak?

Bunun bir kişisel bir de Açık Yazılım Ağı için yanıtı var.

Planlama esnasında naçizane belirli ilkeler etrafında kalmaya özen gösterdim,
ekipte koordinasyon tarafındaki diğer arkadaşlarımın da hemfikir olduğunu
gördükçe mutlu oldum:

- Yaptığımız çalışmalarda tekrarı önlemek ve odaklanmak,

- Gönüllü olan kişilerin sırtından angarya işleri mümkün olduğunca kaldırmak,

- Yasal zeminde kalmak,

- Gönüllü olarak katkıda bulunan insanlarda kötü bir his oluşturmamak,

- Taşıyabileceğimiz kadar yük kaldırmak,

- Açık kaynak olarak devam etmek,

- STKlara/yetkili kurumlara/akredite kurumlara destek vermeye devam etmek,

Plan listesi olarak somutlaştırmak gerekirse,

**Eser Özvataf olarak kişisel kısım:**

- Açık Yazılım Ağı hızlı başladı. Kendi adıma, bunun bundan sonra daha
  sürdürülebilir olması için çalışacağım

- Ekipten ikna edebildiğim arkadaşları ara ara
  [YouTube kanalım](https://eser.live)’a çağırıp Açık Yazılım Ağı için online
  hackathon’lar düzenleyeceğim

- Bazı kurumlardan destek alarak 2-3 ayda bir fiziksel Açık Yazılım Ağı
  Hackathonları düzenleyip odak projeler çıkartarak, yalnızca afet anlarında
  değil, sürekli bir açık kaynak bilinci oluşmasını sağlayacağım
  (başarabilirsem)

- Süreçte olduğu gibi, bundan sonra da kurumlar ile irtibatta olup ortak
  projeler yapmak, açık kaynak projelere katkıda sağlamaya çalışmak, içlerinde
  “Açık Kaynak Geliştirici” istihdam etmek için onları ikna etmeye çalışacağım

- Apache Foundation, CloudNative Computing Foundation, Free Software Foundation,
  Linux Foundation, Mozilla Foundation v.b. ne ise Türkiye’ye de bir yazılım
  vakfı kazandırmaya çalışacağım (işin hukuki/resmi tarafına değil fikir ve çatı
  kısmını inşa etmeye odaklı durumdayım)

- Bu ve benzeri durumlar için gelecekte kişisel olmayan verileri paylaşmak için
  nasıl bir veri platformu inşa ederiz’i planlamaya çalışacağım

- Açık Yazılım Ağı’nın mevcut çözümlerinin “paketlenmesi” kısmında çalışmalarım
  olacak

- Açık Yazılım Ağı’nın portföyündeki yazılımları tanıtıp bunları kullanıcılara
  ulaştıracağım

- Açık Yazılım Ağı’nın çözebileceği potansiyel sorunlar, durumları keşfe
  çıkacağım

**Açık Yazılım Ağı kısmı:**

- Hali hazırda olan
  [code of conduct](https://10forward.io/agreements/code-of-conduct/)’ımızın
  üzerine bir inclusion policy yayınlayacağız (her dilde) (bkz:
  [https://github.com/mozilla/inclusion](https://github.com/mozilla/inclusion))

- Süreç boyunca kendi mesai saatlerinde dahi gece/gündüz çalışan arkadaşlarla
  iletişimi senkron iletişimden asenkron iletişime döndürülecek

- Çok tekrar eden süreçler için akış yapıları kurgulayıp formlarla ilerlenecek

  - Yeni proje formları

  - Yeni katılımcı formları

  - Kaynak ihtiyaçları

  - İş talepleri

  - Ve benzerleri...

- Sürecin hukuki kısmı ile ilgili gereken yükümlülükler yerine getirilecek
  (kişisel verileri yok etme v.b.). Bunun için bir deklarasyon çalışmamız
  gelecek hafta başına kadar yayınlanacak

- Projeler belirli portfolyolara oturtulacak. Bunun için web sitesi üzerinde
  projelerimizin olduğu çalışmamız gelecek hafta başına kadar yayınlanacak

- Projeler tek bir çözüm gibi paketlenecek, bir araya getirilecek

- Paketlenen çözüm white-label ürünler olarak kurulabilir hale getirilecek (bir
  wordpress benzeri blog kurar gibi düşünün)

- Discord çatımızın altında yer alan kitle ile bağlantıda kalmanın yollarını
  arayacağız

**Kısacası** aslında “biz tam gaz devam edeceğiz, KVKK, GPDR başta olmak üzere
yasalara uyumlu bir şekilde, STKlara/yetkili kurumlara/akredite kurumlara
çözümler sağlayacağız. projeler konusunda da odak sağlamak için birkaç gün
içinde yeni bir yönelge yayınlayacağız”.

## Bize nasıl destek olunabilir?

Öncelikle gönüllü başladık ve servis sağlayıcı kurumların bizlere açtığı bazı
hesaplar dışında hiçbir bağış, teşvik, fon v.b. ile ilerlemedik. Bir yazılım
topluluğunun ötesinde davranışlarımız olmadı.

Dolayısıyla bir yazılım topluluğunun temel ihtiyaçları olan bilinirlik, çevre
oluşturma, katılımlarla büyüme ihtiyaçlarımız halen devam ediyor. Özellikle
projede yer aldığını bildiğiniz insanları takip ederseniz zaten destek olmuş
olacaksınız. İlk olarak Twitter, ikinci aşamada Instagram bu süreçte en net
kitlesel aktarımın olduğu platformlar oldu. Bu kanallardan açık yazılım ağı
gönüllülerini ve açık yazılım ağının kendisini takip edebilirsiniz. Ben bilhassa
bu süreçte [kendi Twitter hesabım](https://twitter.com/eser)da yazmaktansa diğer
arkadaşlarımızın tweetlerini dolaşıma sokmaya odaklandım (zaten birçoğumuzun
tweet yazacak zamanımız olmadı).

Farkındaysanız destek için sosyal medyaya yönlendirdim. Çünkü somut bir plan
çıktığında Açık Yazılım Ağı’nın eylemlerini duyuracak bir kanal oluşturmaya
ihtiyacımız var. O nedenle iletişimde kalalım!

## İletişim

İlk günden itibaren sunucudaki moderasyonu tek kişi düzenlemeye çalıştığım an
bile aklımda en kritik işlerden biri “iç iletişim” idi. Discord’un “sahne”
özelliğini oluşturmak, ve o kısmı Furkan’a devrettiğimiz için periyodik
güncellemeler verme konusunda bayağı baskılamış olabilirim. Daha sonra bu
Discord içerisinde bir ritüel haline geldi.

Şimdi iletişimi bir halka daha şeffaf, erişilebilir ve dışa dönük yapma zamanı.

Yukarıda da bahsettim, buradaki eforlara dair hem geri bildirimler almak hem de
daha net konuları aktarmak, hem işin başındaki insanlarla iletişimde bulunmak
için kendi adıma Açık Yazılım Ağı ile ilgilendiğim süreçte
[https://eser.live](https://eser.live) adresinde yayında olmaya çalışacağım.
Eminim bu yayınlarda topluluktan birçok kişi bana katılacaktır, onların
taraflarını dinleyip bir yandan da neler yapabiliriz birlikte konuşmuş olacağız.

## Kapanış ve kişisel öngörülerim

Korkunç bir afet yaşadık. Birçoğumuz için bunun travmatik sonuçları olacak.
Atlattığımız her badirede olduğu gibi biz her ne kadar eski hayatlarımıza dönmek
istersek isteyelim, bir şeyler değişecek veya kendimizi aynı noktada bulmamak
için değişmek zorunda olacak.

Açık kaynak ile tanıştığım yıllarda, okulunu henüz bitirmemiş alaylı bir
yazılımcı olduğum ve çalışmaya devam ettiğim için “uzaktan çalışma”yı da pratik
ediyordum. Pandemi tüm dünyada birçok hasar bıraktı, uzaktan çalışma artık
garipsenmeyen bir pratik olarak hayatımıza girdi.

Açık Yazılım Ağı da kötü bir vesileyle olsa da bizim açık kaynak’ın önemini ve
birlikte öreceğimiz sistemleri tanımak için değişimde rolü olan bir çark haline
gelebilir. Yeter ki biz coğrafya olarak aynı kaderi tekrar yaşamamak için bir
şeyler yapmak isteyelim.

## Sıklıkla sorulan sorular

**Birlikte yapılan ürünlerin kişisel PR malzemesi olduğunu düşünüyor musunuz?**

Bunun hatalı bir soru olduğunu düşünüyorum. Açık kaynak geliştirmede herkes
kendi katkısını ortaya koyarken bundan kıvanç duyabilir, çalışmanın spesifik bir
bölümüne vurguda bulunup kendini veya yakın arkadaşlarının çabasına daha fazla
değer atfedebilir. Bana kalırsa projenin tamamına görünürlük sağladığı için
bunun çok da büyük bir sakıncası yok. Siz de projelerin doğru tanıtılmadığı,
eksik kaldığı kısımlar için tanıtım yazabilir, README’leri güncelleyebilir,
sosyal medyada bahsedebilirsiniz. Ben elimden gelince kendi kaynaklarımla
bunları desteklemeye çalışıyorum.

Herkes bu çalışmalar esnasında tanıştığı kişileri etiketleme, onlara
minnettarlık gösterme anlamında serbest olursa bu daha çok olumlu bir paylaşım
doğurur. Elbette buradaki 24000’e yakın kişiyi tek tweet’de etiketleyemeyiz ama
herkesin katkısı burada, aynı zeminde olarak zaten eşit derecede kıymetli.

**Neden formlarla iletişim yapıyoruz?**

İlk 1 hafta boyunca birçok gönüllü ellerinden geldiğince insan üstü eforlarla
bir çalıştılar. Ancak ikinci hafta itibariyla ulusal yas’ın da bitmesiyle hem
daha az zaman ayırabilir hem de nöbetleşe çalışır bir düzene dönmek zorunda
kaldılar. Bu anlamda taleplerin, “cloud ekibi”, “ML ekibi” gibi nöbetleşe
bakabileceği ve sistem üzerinden kontrol edebileceği ve kaybolmayacağı asenkron
iletişimle ilerlenebileceği bir yapıya geçmek istiyoruz. Bu işin
sürdürülebilirliği için şart görünüyor.

**Neden bazı projeler iptal ediliyor?**

Biz istiyoruz ki, bir sonraki afette hazırlıksız yakalanmayalım. Hatta
Kahramanmaraş depreminde üzerimize düşen vazifeyi tamamladıktan sonra olası bir
sonraki afette yine tüm her şeye 0’dan başlamayalım.

Coğrafyamızı koruduktan sonra, gerekirse dünyanın diğer taraflarında da
kullanılabilecek yangın, tsunami, fırtına, yangın v.b. felaketlerde de insanlığa
yardımcı olabileceğimiz bir açık kaynaklı, Wordpress gibi kurulabilen, bize
bağımlı olmayan bir sistem geliştirelim.

Bunu yaparken de odaklanmak ve birbirinin benzeri işleri yapan proje ekiplerini
birleştirmek, devamında da insan ve hesaplama/computing kaynaklarımızı efektif
kullanmamız gerekmekte.

Aslında burada bir iptal değil, yazılım vakıflarının da yaptığı gibi projeleri
basamak sistemiyle sıralama söz konusu olacak.

**Biz form doldurduk ama geri dönüş olmadı. Neden önemsemediniz?**

Önemsedik, ancak 4 günde herhangi bir kurumsal firmanın bile kolay kolay
altından kalkamayacağı bir organizasyon ve görev dağılımı oluşturduk. İnanın bu
tarihte birkaç yerde “case study” olarak anılacak bir çalışma ve gönüllüler
olarak hepimiz bunun bir parçası olduk.

Ancak bir süre sonra STKlarla/yetkili kurumlarla/akredite kurumlarla
irtibatlanarak mevcut talepleri sahaya yansıtmaya odaklanmamız gerekiyordu ve
bir noktada proje sayısını limitleyerek yeni proje alamaz hale geldik. Özellikle
altyapı ekiplerimiz orada güvenlik kontrolleriyle birlikte insanüstü bir eforla
çalıştılar. İnanın zaman sınırına takıldık. Siz açık yazılım ağı çatısı altında
projelerinizi sürdürdüğünüz sürece, ben şahsen birkaç hafta içinde gelip sizinle
tanışmak ve konuşmak için elimden geleni yapacağım. Eminim diğer arkadaşlarım da
dönüş gerçekleştirecektir. Eğer ulaşamazsam twitter üzerinden de
erişebilirsiniz. Halen tüm mesajlara dönemedim. Kusura bakmayın.

## Bazı bağlantılar

Entegre yapılar:

- [https://www.mynet.com/deprem-yardim-agi](https://www.mynet.com/deprem-yardim-agi)

- [https://crypto.ahbap.org/en](https://crypto.ahbap.org/en)

Videolar:

- **Deprem Dayanışma YouTube Ortak Yayını**\\
  (linki bulamadım, güncelleyeceğim)

- **Webtekno röportajımız**\\
  [https://www.youtube.com/watch?v=2rmDUQYQc3k](https://www.youtube.com/watch?v=2rmDUQYQc3k)

Bizden bahseden makaleler:

- **WIRED**\\
  [https://www.wired.com/story/tech-volunteers-rush-to-save-turkeys-earthquake-survivors/](https://www.wired.com/story/tech-volunteers-rush-to-save-turkeys-earthquake-survivors/)

- **TIME**\\
  [https://time.com/6254500/turkey-earthquake-twitter-musk-rescue/](https://time.com/6254500/turkey-earthquake-twitter-musk-rescue/)

- **Euronews**\\
  [https://www.euronews.com/next/2023/02/10/how-twitter-helped-find-survivors-trapped-beneath-rubble-after-turkeys-earthquakes](https://www.euronews.com/next/2023/02/10/how-twitter-helped-find-survivors-trapped-beneath-rubble-after-turkeys-earthquakes)

- **Google Türkiye Resmi Blog**\\
  [https://turkiye.googleblog.com/2023/02/yasanan-deprem-felaketinin-ardndan.html?m=1](https://turkiye.googleblog.com/2023/02/yasanan-deprem-felaketinin-ardndan.html?m=1)

- **Milliyet**\\
  [https://www.milliyet.com.tr/yazarlar/hanife-bas/yazilim-ordusu-seferber-oldu-6902301](https://www.milliyet.com.tr/yazarlar/hanife-bas/yazilim-ordusu-seferber-oldu-6902301)

- [https://www.ekonomim.com/kose-yazisi/girisimci-z-kusagi-tarih-yaziyor/682924](https://www.ekonomim.com/kose-yazisi/girisimci-z-kusagi-tarih-yaziyor/682924)
', 'tr          ', '01J6ZJ7RMK5Z7ETSNTHJYAVZ3Z', 'AYA''nın oluşum hikayesini kaleme alıp, bundan sonrası için hedeflerini anlatmaya çalıştım.', 'Açık Yazılım Ağı');
insert into "story_tx" ("content", "locale_code", "story_id", "summary", "title") values ('Bugün yazılım sektöründe çalışan tanıdıklarla bir araya geldiğim mini
buluşmalardan birindeydim :) Oldukça keyifli konular konuşma fırsatımız oldu.

Bu konulardan bir tanesi de “Mentorluk” üzerineydi. Mentorluk, daha önce sevgili
Fatih Kadir Akın ile de oldukça üzerine kafa patlattığımız bir kavramdı. Çünkü
bu sıfat bizce “birinin hayattaki akıl hocası olmak” gibi altında büyük bir
anlama sahip olduğu için, ancak mentorun hayatta karşısındakinden çok daha
deneyimli, karşı tarafın zihnini açan ve alanında otorite olduğu kabul edilmiş
birisi olduğunda kullanılabilecek bir ağırlıktaydı.

Daha sonra kurumlar, departmanlar, topluluklar “mentor-mentee” programları
başlatıp, sektörde birçok uzman kişinin kendisine “mentor” diyemediği bir
ortamda bu kavramı neredeyse bir fast-food ürünü gibi kullanmaya başladılar.

Her ne kadar yukarıdaki ifadelerim “olumsuz bir yargı” içeriyor gibi tınlasa da,
aksine bunun işe yaradığını, birçok kişiye yardımcı olduğunu da gözlemledim.

Ancak içeriğe baktığımızda, buradaki iletişimin mentorluk gibi “hayat görüşü”
bağlamında değil de, daha çok birçok kişinin bilhassa yazılımın ilk dönemlerinde
konulara nasıl yaklaşması gerektiğini aktaracak kapsamda bir “koçluk” olduğunu
fark ettim.

Kavramsal olarak iki kavramın arasındaki farkı Wikipedia alıntıları ile şöyle
aktarmak isterim:

[Koçluk](https://tr.wikipedia.org/wiki/Ko%C3%A7luk):

> Bazen koçluk, iki kişi arasında gayri resmi bir ilişki anlamına gelebilir.
> Biri diğerinden daha fazla deneyime ve uzmanlığa sahip olabilir, ikincisi
> öğrendikçe tavsiye ve rehberlik sunar. Fakat koçluk, daha genel hedeflerin
> veya genel gelişimin aksine, belirli görevlere veya hedeflere odaklanarak
> mentorlukdan farklıdır.

[Mentorluk](https://tr.wikipedia.org/wiki/Mentorluk):

> Geniş deneyime sahip biriyle öğrenmek isteyen biri arasında bir öğrenme ve
> gelişme ortaklığıdır. Bir uzmanla etkileşime geçmek, kültürel araçlarla ve
> içinde uzmanlık kazanmak için de gerekli olabilir. Mentorluk tecrübesi ve
> ilişki yapısı “psikolojik sosyal destek, kariyer rehberliği, rol modelleme ve
> protégé’lerin ve mentorlerin dahil olduğu mentorluk ilişkilerinde meydana
> gelen iletişimi” etkilemektedir.

Kısacası “koçluk” birinin “yazılıma başlaması”, “kariyer hedeflerini
belirlemesi” gibi sonuç odaklı bir pragmatik ilişki iken; “mentorluk” biraz daha
kültürel ve psikolojik gelişime kadar kişinin kendisini yeniden
tanımlayabilmesine kadar uzanacak bir rol modellik içermekte.

Yazının ana noktası oldukça kısa olduğu için bu kısalıkta tutmanın faydalı
olacağını düşünüyorum, kalın sağlıcakla 👋
', 'tr          ', '01J6ZJ7Q8MYTANHEV0THKJ2XVN', '"Mentorluk" sıfatı, sektörde son yıllarda hızlıca tüketilen bir kavrama dönüşmüş durumda. Bu kavramların kökenlerine inerek iki kavramın farkına dikkat çekmek istedim.', 'Mentorluk ve Koçluk');
insert into "story_tx" ("content", "locale_code", "story_id", "summary", "title") values (E'Başlıklarda bir “yolculuğumuz” teması tutturdum gidiyorum, farkındayım. Çünkü
“süreç” kelimesinin “katlanılması gereken bir dönem” gibi bir alt okuması,
yolculuk kelimesinin ise içinde öğrenimi ve benimsemeyi barındırdığını
düşünüyorum. Bu nedenle bu dönem yolculuk kelimesini çokça kullanmış olabilirim.

Her ne kadar bana Twitter’da rastlayanlar bilfiil kod yazdığımı düşünüyor
olsalar da, son 10 yıldır profesyonel olarak yazılımın daha çok “soft skill”i
yüksek taraflarında çalışıyorum. Bunu yaparken de mevcut yazılım
organizasyonlarını, ekiplerini bir başlangıç noktasından hedeflenen “daha olgun”
bir noktaya taşımak için belirli stratejik / taktik müdahalelerde bulunmam
gerekebiliyor.

“Durumsal liderlik”i içselleştirdiyseniz, bu müdahalelerin ekibin yapısı, takım
kimliği, farkındalığı v.b. etkenlere göre değişkenlik arz edebileceğini; “dünden
daha iyi” bir noktaya giderken bazen otokratik, bazen demokratik, bazen laissez
faire karar alan bir liderlik sergilemeniz gerektiğini bilirsiniz. Ben de bu
hareket tarzını benimsiyorum.

Ancak “karar alma” bazen işin “en basit” noktası olabiliyor. Bu nedenle bu
yazının odak noktası daha çok “kararı uygulamak”la ilgili olacak.

---

Medyanın ağzına sakız ettiği Z jenerasyonu, benim de mensup olduğum öncül’ü
millenial jenerasyonuna göre birkaç farklı davranış kümesine sahip -ki bu doğal
bir durum-.

“Benim naçizane deneyimlerimce” millenial’lar kendi öğrenme yolculuklarında
kendi öncüllerinden devraldıkları “usta’yı izlemek” veya kaynak yoksunluğunda
son çare olarak “inspect & adapt” (gözlemle ve adapte ol) yaklaşımlarına daha
yatkınken, Z jenerasyonunda fikirlerle zaman kaybetmeden önce o fikirlerin işe
yarayacağına dair “kuşkucu” bir yaklaşım benimseniyor. Elbette tam tersi şekilde
milleniallar içinde kuşkucular, Z jenerasyonu içinde takipçiler de var, ancak
ben gözlemlerim temelinde çoğunluğun davranışını ele alıyorum.

Bu kuşkucu tutumla karşılaştığımızda, “tatmin” etmek için iyi bir hitabet,
örnekler ve rasyonelite sunmak neredeyse şart. Ancak takdir edersiniz ki, bunun
için her zaman “zaman bulmak” büyük bir sorun ve zaman olsa dahi “operasyonel”
konuların sürekliliği gerekirken antik yunan dönemi filozofları gibi
argümanlarla uzun uzun tartışmak pek pragmatik/yararcı olmuyor.

---

Karşımızdaki sınav bu olduğunda; scrum çerçevesi, kuşkuculuğun kanıt-tabanlı
yapısını kucaklayarak empiricism/deneyselcilik-tecrübecilik’i öne sürüyor. Ve
özetle diyor ki; dene, bir şeyler yanlış gidiyorsa yeni bir şey denersin. Yani
kaizen’den aldığı esinlenme ile bizi bir “inspect & adapt” yöntemine sürüklüyor.
İşin çok felsefesine girmek istemesem de, “bilemeyeceğin” bir kavrama karşı
belirsizliği yok etmenin en iyi silahın “bilinmezliği ortadan kaldırmanın”
olduğu düşüncesi oldukça baskın Scrum’da. Ve ben de bunu günlük hayatımda bir
ilke olarak kullanıyorum. Ancak “bir şey öğrenme yolculuğu”nda bu prensip bence
biraz eksik kalıyor.

Naçizane, empiricism’in uzun soluklu bir öğrenme yolculuğunda mindset/düşünce
yapısı olarak bizi yanlış yönlendirebileceğini düşünüyorum. Nedeni ise şu:
empiricism’de “deney süresi”ne dair bir yönlendirmeye ulaşamıyoruz. Elbette bu
yaklaşımın bir zaafı değil, yalnızca bana “öğrenmeye” uygulamaya uygun
olmadığını düşündürüyor. Örneklendireyim: hiçbir müzik bilginizin olmadığını
varsayalım, elinize gitar’ı aldıktan sonra gitar’ın sizin çalmanız için uygun
bir enstrüman olmadığı yargısına ne zaman varabilirsiniz?

---

Bu noktada yine uzak doğudan ithal “Shu Ha Ri”den bahsetmenin sırası geldi.
Kapak fotoğrafında biraz ipucunu verdiğim üzere Shu Ha Ri’nin Aikido gibi savaş
sanatları ile bir bağlantısı var. Bizler için bir disiplini ustasından öğrenip
ustalaşma yolculuğunu tanımlar.

**Shu**: Usta’nın sana öğrettiklerini değiştirmeden uygula\\
**Ha**: Uyguladıklarının arkasındaki fikri ve felsefeyi anlamaya ve
içelleştirmeye çalış\\
**Ri**: Ezbere yaptıklarından kurtul, ilgili fikir ve felsefe çerçevesinde kendi
yolunu bul

Ben elbette bunu çekip çıkarmış bir sektör dahisi değilim, Alistair Cockburn ve
Martin Fowler gibi isimlerden beslenerek sizlere bu yazıyı ulaştırabiliyorum.

Ancak tabii bunun üzerine koyabileceğim kendi deneyimlerim de var. Yukarıda
bahsettiğim gibi “dünden daha iyi bir noktaya ulaşma” hedefiyle yola çıkılan
ekiplerde, kuşkucuların her zaman “bu bizde çalışmaz”, “biz bunu denedik
olmadı”, “ya bunu eleştiren şöyle şöyle bir blog yazısı var”, “kural koymak
otonomiyi öldürüyor” v.b. argümanlarla çokça baştan “yapmayalım” mindset/düşünce
yapısı içinde olması durumunu yaşadım.

Bilhassa sektörde “best practice” olarak benimsenen konuların neden “genelde
başarılı” olup sanki diğer ekiplerin “özel durumları” yokmuşçasına yalnızca o
ekip için “özel durumlar nedeniyle uygulanamayacağı” fikrinin çok geçersiz
olması bir yana dursun; “peki bu yöntemi kullanmayalım, ama bir yöntem
kullanmalıyız değil mi? yerine ne kullanalım?” sorusunun da yanıtsız kaldığını,
dolayısıyla kuşkucuların aslında “literatüre hakim olmadan” peşinen olumsuz
yaklaştığı maalesef genel-geçer bir durum.

Başka bir örnek de; kavramları kendi silahıyla vurduğunu düşünen kimselerden
gelsin. Yazılım alanında senior seviyede olsalar da çok arkadaşımdan “madem
Scrum otonom takımlar destekliyor o zaman neden daily yapmak zorundayız” gibi
sorular çok aldım. Her seferinde “aslında zorunda değilsin ama bunu tartacak ve
doğaçlayacak olgunlukta bir ekipten mi bahsediyoruz?...” gibi sorgulamalara
girsem de kendi bildiklerini okuduklarını; insanların genelde merdivenleri adım
adım çıkmak yerine 5-6 basamak birden çıkmaya çalışırken kayıp düştüklerine
şahit oldum.

Peki Shu Ha Ri felsefi olarak ne diyor?

Aslında diyor ki: “sürecin mantığını yargılamayı” ve “sonuç çıkacak mı diye
kaygılanmayı” geride bırak, beyninin şalterini indir. Önündeki kuralları büyük
bir disiplin ile uygula. Ardından disiplinin hayatına getirdiği iyi ve kötü
deneyimleri gözlemle.

Bu deneyimleri yaşamanı sağlayan yöntemler aslında bu kuralların en başında neyi
sağlamak, ve bir önceki çalışma pratiklerinde neyin eksik olduğunu da sana
gösterecek. Hali hazırda hayatında olmayan bir etken için önceden hesaplama
yapmak çok da mümkün değil.

Ancak uyguladığın disiplini tam anlamıyla yapabilir ve sonuçlarını toplayabilir
hale geldiğinde artık kuralları yıkmaya hazır olursun...

Yine gitar örneğine geri dönmemiz gerekirse, parmak egzersizlerini yapmadan,
enstrümanı nasıl tutman gerektiğini bilmeden, doğru akorları basmadan kalkıp
diğer virtüözlerden farklılaşacak bir “keşif”te bulunman hiç olası değil...

---

90’larda büyüyen neslin “ustaya” daha fazla saygı duyması, ancak cep
telefonlarının icadıyla hızlanan yaşamımız, hayat kaygılarımız, internet’in
yaygın kullanımı bilginin kişilerin tekelinde olmaması gibi etkenlerin tümü
toplanarak bize ve bizden sonraki nesillere “aceleci” bir dünyayı empoze etti.
Bu nedenle “Shu Ha Ri”nin önerdiği ustalaşma yolu ne kadar benimsenir bilmiyorum
ama pratikleri içselleştirmek konusunda bildiğim en doğru mindset/düşünce
yapısı.

Kalın sağlıcakla,
', 'tr          ', '01J6ZJ7RSFJV6VBYPT467BKFQ0', 'Mühendislik odaklı "olgunlaşma yolculuğu" üzerine çalışmalarımdan ve son dönem karşılaştığım durumlar ertesinde biraz Shu Ha Ri kavramını yazıya almak istedim. Bir disiplini gerçekten benimsemenin yollarından biri, onu önce sorgulamadan uygulayıp (Shu), ardından anlamlandırıp (Ha), sonunda da özgünleştirerek kendi yolunu çizmektir (Ri). Bu yazıda bunu anlatıyorum.', 'Shu Ha Ri, Empiricism, Kaizen, Inspect & Adapt - Disiplin Öğrenme Yolculuğumuz');
insert into "story_tx" ("content", "locale_code", "story_id", "summary", "title") values ('Konuya oldukça kişisel bir bakış açısından girip, daha sonra toparlayacağım.
Şimdiden girişteki uzun hikaye kısmı için kusuruma bakmayın :)

![](https://cdn.hashnode.com/res/hashnode/image/upload/v1659040683473/nb6DGQfNf.png)

### **Açık Kaynak Yolculuğum**

İstisnasız davet edildiğim her konuşmada ilk 2–3 dakikada kendimi tanıtırken
kullandığım bir cümledir, “açık kaynak ekolünden geliyorum”. Pek irdelemem,
detayına girmem ancak yalnızca bir cümle olsa dahi geçirmekte de fayda görürüm.
Bunun birçok nedeni var, ancak en kişisel olanından bahsedeyim; Kağıt üzerinde
kurumsal firmalarda çalışan bir beyaz yakalı çalışanım. Bu da yetmezmiş gibi
günlük mesaimde gerekmedikçe ellerimi çok az kirletiyorum. Dışarıdan yazılım
üretimi -istese de- şirketinin orta vadeli hedeflerinden _daha öte_ bir değere
dönüştüremeyen ve 5 yıl öncesine kadar açık kaynak bir koda güvenilemeyeceğinin
iddiasında olan biri gibi değerlendirilmem için tüm şartları sağlıyorum.

Oysa ki 90’ların ortasında BBS’ler döneminde yazılımla uğraşmaya başlamış ve
yazılım geliştirmeyi topluluktan öğrenmiş ve topluluğun yararını çokça görmüş
biri olarak işin “sosyal” boyutunu atlamak ve kocaman bir yazılım ekosistemini
tek taraflı kullanmak benim için bir seçenek dahi olamazdı. (bkz:
[90’larda programlama öğrenmek](../yazilim-araclari/20160221-90larda-programlama-ogrenmek)
yazım)

Disket/CD değişimleri, internette kişisel web sitelerinde yayınladığımız kodlar,
phpclasses.org ve sourceforge maceralarımı hızlıca atlamam gerekirse; bugünkü
manada açık kaynak ekosistemine ilk katkımı 2011’de tanışıklığımın olmadığı, bir
çıkar ilişkisi beslemediğim insanların GitHub repository’lerine katkı yaparak
gerçekleştirdim.

O günden bu yana da yazılımla ilgili biriktirdiğim tüm bilgi ve deneyim
birikimimi -gündüz yapmış olduğum mesaiden komple bağımsız bir biçimde- kendi
kişisel zamanımdan arttırarak açık kaynak gönüllüsü olarak çalışıyorum. Hızlıca
bir baktığımda GitHub’da “Repositories” bölümünde 189 projenin hesabıma tanımlı
olduğunu söylüyor, bunların yalnızca 27’si fork. Ve ben halen daha üretken
olabileceğimi, henüz elimde bir şey olmadığını düşünüyorum. Çünkü yaptığım ufak
tefek çalışmalarda bir bütünlük eksikliği var. Hepsinin bir amaca hizmet
etmesini amaçlıyorum.

Diğer taraftan genel bir dönüşüm olduğunun, “çok büyük katkısı olmayan” biri
olarak benim koyduğum tuğlaların bile nihayetinde kültürel anlamda bir işlevi
olacağından eminim. Çünkü en basitinden bir önceki paragrafın başında değindiğim
“açık kaynak projelerin kişisel zamandan arttırarak yapılan bir aktivite olması”
konusunda 5–6 yıl sonra bir şeyler daha farklı olursa, açık kaynak ekosistemine
dahil olan daha fazla kuruluş olursa, bunun bugünün açık kaynak gönüllülerinin
başarısı olduğundan emin olacağım.

### **github.com/acikkaynak**

Bu fikirlere sahip biri olarak yaklaşık 5 yıl önce github.com/acikkaynak
adresini yayına aldım ve çağrıda bulundum. Bugünkü acik-kaynak.org’nun “MVP”
hali olan bu çalışma, insanların bir GitHub repository’sine kendi projelerini
eklemesi ve dolayısıyla topluluğun geri kalanına da “benim de bir açık kaynak
çalışmam var, katkı sağlayabilirsiniz” mesajı bırakabilmesi anlamına geliyordu.

Elbette ihtiyaçlar ihtiyaçları doğruyor, hemen ardından bu eklenen projeleri
JSON formatına taşıdık ve projelerin listelendiği HTML otomatik oluşmaya
başladı. Bir sonraki adım ise artık acikkaynak.info (bugün artık bu adresi
kullanmıyoruz) adresini üzerinden insanların nasıl açık kaynak geliştirmeye
başlayabileceklerini anlatmak ve GitHub’da yer alan açık kaynak proje listesini
bir web sitesi arabirimi üzerinden görmelerini sağlamak oldu.

github.com/acikkaynak/acikkaynak adresi üzerinde 95 açık kaynak proje
listeleniyor. Ve yalnızca Slack üzerinden 391 açık kaynak yazılım geliştiriciye
ulaşıyoruz.

Daha iyi yapabileceğimiz şeyler vardı. Birçok kişi ile tanışma, muhabbet etme ve
projelerinde ufak tefek konulara yardım etme şansı buldum bu süreçte. Bazı
kimselerin projelerini fark edip pull request ile düzenlemeler yapma fırsatı
buldum. Ama birlikte bir proje planlandı veya açık kaynak yeni bir projenin
temelleri bu platformda atıldı mı? Maalesef hayır. Evet ama gördüğüm şu ki:
potansiyelimiz var.

### **acik-kaynak.org**

Bugün [acik-kaynak.org](https://acik-kaynak.org/) adresi üzerinden yeni bir
sürüm üzerinde çalışıyoruz. “Açık kaynak’ı açık kaynak yapıyoruz” diye
özetlediğim bir çalışma gerçekleşiyor.

Agah Gürer’in yaptığı tasarımı Figma’ya aktardık, şu anda herkes ulaşıp
görüntüleyebiliyor.
[Açık Kaynak 1.0 Yol Haritası](https://github.com/acikkaynak/acikkaynak-website/issues/30)
başlığından tüm geliştirme detaylarına ve aktif sprint hedeflerine
ulaşabiliyorsunuz. Okan Davut, Samet Aylak, Mehmet Ali Peker, Berat Bozkurt,
Tuğsan Ünlü, Numan Can İpekçi ve Kaan Gökdemir koda pull requestler aracılığı
ile katkıda bulundular. Arda Kılıçdağı, Emre Yılmaz, Hatice Edis, Emir
Karşıyakalı, Selçuk Ermaya, Selen Gora ve Samet Aydemir’in ise GitHub
Issue’larında izleri var. Aliyar Fırat sürüm planlamada ve tasarımda birçok
fikir getirdi.

Üstelik bu yazının yayınlandığı tarihten 1 hafta sonra da projeyi tamamlamak
gibi bir hedef var.

Ardından sürümlerle gideceğiz ve backlog’umuz şimdiden çok dolu. Bazı fikirler
uçuk kaçık olsa da;

- Açık Kaynak geliştirme hakkında 0’dan başlayan birini donanımlı hale getirecek
  bir içerik sağlayan,
- Kurumlar nasıl açık kaynağa geçiş yapabilir, bu konuda içerik sağlayan ve
  çözümler üreten,
- İnsanların GitHub hesapları ile girip erişime izni oldukları repoları
  işaretleyerek sisteme aktarabildikleri,
- Erişebildikleri repolar hakkında haber girebildikleri, -dolayısıyla başka
  kullanıcıların da gördüğü sosyal bir mecra haline geliyor-
- Sistemdeki projelerin künye sayfalarının olduğu, böylece insanların projeleri
  beğenip, sosyal medyada paylaşabildikleri,
- Sistemdeki kişilerin istatistik bilgilerinin yer aldığı,
- Teknolojilere, kategorilere hatta etiketlere göre aranabilen projeler
  sayfasına sahip,
- Sponsorluğu teşvik eden, sponsorların reklamını yapan,
- Projeler için girilen haberler, istatistikler ve yeni eklenen projelerle
  haftalık açık kaynak bülteni yayınlayan,

yeni sürümler için de yol haritaları oluşmaya başladı. Genel anlamda bu projenin
amacı Türkiye’deki ve Türkçe konuşan geliştiricilerin oluşturduğu açık kaynak
ekosisteminin lokomotifi olmak, insanları açık kaynak geliştirmeye motive etmek,
farkındalık yaratmak ve engelleri ortadan kaldırmak. Dolayısıyla ekosistem neye
ihtiyaç duyuyorsa o noktada görev almaya devam edeceğiz.

### **Diğer Açık Kaynak Projelerim**

Tekrar kendime dönecek olursam bugünlerde acik-kaynak.org’un geliştirmesine
odaklansam da diğer açık kaynak projelerime de dönem dönem ağırlık veriyorum.
Bunları daha önce tanıtacak fırsatım olmamış bu nedenle belli başlı olanlardan
bazıları hakkında yazarak bu eksiği kapatmak istiyorum.

Eğer ilgilenirseniz dilediğiniz projeye sizleri de katılımcı olarak
ekleyebilirim. Hatta telegram üzerinden benimle iletişime geçerseniz konuyu
projelendirir, ortak hedefler koyar ve projeler üzerinde birlikte çalışabiliriz.

[**darty**](https://github.com/eser/darty): (JavaScript) Darty en temel anlamda
bugünün modern web geliştiricisinin ihtiyaçlarını karşılayan bir geliştirme
aracı. TypeScript compiler, Jest ve webpack kullanarak sıfır konfigurasyon ile
SSR, Module Bundling, CSS Modules, Hot Reloading, Containerization ve TDD-Ready
Environment sağlıyor.

Belirli darty şablonları var, “darty-react-app”i örnek verelim: içerisinde react
ve semantic ui react ile birlikte geliyor. Vue, Svelte, v.b. şablonlar
oluşturmak için ortada hiçbir engel yok. Bu şablonlardan birini seçip “yarn
create darty-app react webapp” gibi bir komutla çalışmaya hazır bir tabanla
hemen başlayabiliyorsunuz. Create React App’a göre avantajları/dezavantajları
olan ancak olaya tamamiyle farklı yaklaşan bir araç darty.

[**evangelist**](https://github.com/eser/evangelist): (JavaScript) Fonksiyonel
programlamaya aşinaysanız ve compose, curry, decorate v.b. yapıları fazlaca
kullanıyorsanız tree-shaking uyumlu bu kütüphaneyi sevebilirsiniz. loadash gibi
komple bir alet çantası yerine inline edebileceğiniz sadelikte fonksiyonları
size sunuyor.

[**senior**](https://github.com/eser/senior): (JavaScript) ufak bir uygulamanız
var diyelim ve buna üçüncü parti uzantı (extension) desteği yazmak istiyorsanız.
bu proje npm paketleri aracılığı ile uygulamanızın uzantılarını
listeleyebileceğiniz, indirip kurabileceğiniz ve kaldırabileceğiniz bir arabirim
sunuyor.

[**consultant**](https://github.com/eser/consultant): (JavaScript) belirli bir
parametre diziniz var ama bu parametreleri ister bir JSON girdisinden, ister
komut satırından, isterseniz de kullanıcı ile iletişime geçip GUI üzerinden
almak istiyorsunuz. consultant sizin için bu girdi esnekliğini sağlamaya
çalışıyor.

[**maester**](https://github.com/eser/maester): (JavaScript) loglama ve
exception/error handling işlerini aynı potada eriten, hem node.js hem de
browser’da aynı kod tabanı ile çalışan bir kütüphane.

[**enthusiast**](https://github.com/eser/enthusiast): (JavaScript) fonksiyonel
yaklaşımları kullanan ancak kaynak olarak da streaming arabirimleri kullanan bir
kütüphane.

[**ponyfills**](https://github.com/eser/ponyfills): (JavaScript) çoklukla
kullanılan ancak tüm platformlarda desteklenmeyen belli başlı JavaScript
foksiyonlarının eğer varsa native halini kullanan, yoksa bunların çoğunlukla
mozilla developer network üzerinden alınan polyfill sürümlerini çalıştıran
kütüphane.

[**immunity**](https://github.com/eser/immunity): (JavaScript) immutable veri
yapılarını koruyarak array ve object gibi nesnelerinde işlemler yapmanızı
sağlayan bir kütüphane.

[**servicemanager**](https://github.com/eser/servicemanager): (JavaScript)
JavaScript için inversion of control sağlayan bir kütüphane.

[**cofounder**](https://github.com/eser/cofounder): (JavaScript) Dosya sistemi
ve OS işlevlerini soyutlayarak bunları implementasyondan bağımsız olarak
kütüphanelerin kullanabilmesini sağlayan bir katman.

[**sey**](https://github.com/eser/sey): (JavaScript) Bir web projesinde oldukça
fazla tekrarlanan işleri bir task tool ile değil de, declarative olarak
belirterek çözmenizi sağlayan bir araç.

[**apibone**](https://github.com/eser/apibone): (JavaScript) Slack, Telegram,
Webhook’lar v.s. gibi ortamlardan yapılan çağrıları tekilleştirerek arkada bir
bot mekanizması işletebilmenizi sağlayan, bot’dan alınan yanıtı yine ilgili
katmanın kurallarını işleterek kullanıcıya ileten bir uygulama.

[**tassle**](https://github.com/eser/tassle): (.NET Standard) Konfigurasyon,
Data, Fonksiyonlar, I18N, Loglama, Doğrulama gibi kütüphaneler sağlayan bir .NET
paketi.

[**jsmake**](https://github.com/eser/jsmake): (JavaScript) JavaScript için make,
rake gibi bir task runner. Ancak kullanıcının tanımladığı taskları farklı
kanallardan okuma, uzantı desteği gibi ek özelliklerle donatılmış halde.

Diğer projeler için [https://github.com/eser](https://github.com/eser) adresini
inceleyebilirsiniz.

### **Son Olarak**

Tahminen açık kaynak üzerine çokça çalışmaya, konuşmaya (bkz:
[Kurumlar için Açık Kaynak Pratikleri Sunumum](https://speakerdeck.com/eser/acik-kaynak-pratikleri))
ve içerik üretmeye devam edeceğim. Lütfen ortak bir şeyler yapabileceğimizi
düşünüyorsanız [https://t.me/acikkaynak](https://t.me/acikkaynak) adresinden
Telegram üzerindeki “açık kaynak development” kanalına katılın, tanışalım.
', 'tr          ', '01J6ZJ7RXSVEDNK55KKZZ8VZ2N', 'Açık Kaynak İnisiyatifi''nin açık kaynak kültürünü yaygınlaştırmak, katkıyı kolaylaştırmak ve Türkiye’deki geliştirici ekosistemine sürdürülebilir bir alan açmak hedeflerini anlattım.', 'Açık Kaynak İnisiyatifi');
insert into "story_tx" ("content", "locale_code", "story_id", "summary", "title") values (E'Steve Ballmer zamanında “sözde”, Satya Nadella ile birlikte ise kısa zamanda
“özde” geliştirici dostu olan Microsoft’a açık mektubumdur.

Öncelikle hayatımda profesyonel olarak kazandığım gelirin büyük çoğunluğu
Microsoft teknolojileri üzerinden oldu, 1997 yılından bu yana Microsoft
teknolojilerini kullanarak yazılım geliştiriyorum. “Microsoft yazılımcısı”
olduğum kadar, son 5–6 senelik yükselen trendlerle birlikte GNU/Linux — Özgür
Yazılım cephesinde de dengeli bir şekilde yazılım geliştirmeye başladım. Artık
favori bir platform benimsemek yerine birden fazla alanda üretim yapıyorum. Ve
eminim tek alanla kısıtlı kalmayan, birden fazla platforma hakim olan tek
yazılımcı değilim.

Microsoft şu ana dek kendi dünyası dışındaki araçların çoğunu “port” etti. .NET
oldukça güçlü bir oluşumdu ve eğer eksik olduğu bir cephe oluşursa anında kendi
alternatifini oluşturuyordu. Hibernate yerine NHibernate, GitHub yerine
CodePlex, PIP/NPM yerine NuGet, Make yerine MSBuild, StyleCop, NUnit, Entity
Framework vb. Bunlardan bazıları benim nezdimde sınıfta kalmış, bazıları ise
güçlü birer araç olarak halen varlığını korumaktadır.

Uzun süredir ASP.NET yazarken bile Grunt’dan, NPM’den, PhantomJs’den yardım
alırdım ama bu araçlar Microsoft platformları için üvey evlatlardı. Fakat şu
sıralar ASP.NET MVC, Azure ve son olarak VS 2013 ile bir şeylerin değiştiğini
hissedebiliyorum. Yalnızca Microsoft’un geliştirdiği değil, topluluktan gelen,
topluluğun desteklediği araçlar da kendine .NET ve Microsoft ekosisteminde yer
bulabilmeye başladılar.

Fakat ne yazık ki Windows halen bir geliştirici için çağın gereksinimlerini
karşılayabilen bir platform değil. Yanlış anlaşılma olmasın, Windows inanılmaz
olgun bir sistem, bir çok avantaja sahip fakat “eksik”. Ben bu eksikler
arasından bana göre en kritik olanına, komut satırına odaklanmak istiyorum.

“Madem Microsoft ASP.NET vNext açılımını yapacak kadar olgunlaştı, neden
yazılımcılar için daha fazlası olmasın?” motivasyonu ile Windows 10’un bir milat
olmasını beklerken, ben kendi bulabildiğim eksikleri saymaya çalışacağım.

**Güçlü Komut Satırı:**\\
Öncelikle bana lütfen PowerShell demeyin. PowerShell’e işiniz düşmediği sürece
PowerShell kullanmazsınız, ben şahsen yolunu bile bilmiyorum. İşletim sisteminin
kendi shell’ini esas alıyorum. Windows 10’da ilk kez komut satırı için bir
geliştirme gelecek yıllar sonra, o da copy-paste tuşları. Komut satırında her
şeyden önce bash kadar yetenekli olmalı. Neden “echo \\`ls\\`” gibi bir komut
veremiyorum? Neden aliaslar yok?

**Komut Satırı Renk Kodları:**\\
16’lık ANSI renk kodları ben Quick Basic kullanırken yeterli idi. Fakat
88/256’ın şu anda geçerli olması gereken standart olduğuna inanıyorum.

**GNU Araçları:**\\
Şu anda GNU araçlarını haricen Windows’a kurup çalıştırabiliyorum. Oysa ki
bunlar işletim sistemi ile birlikte gelebilir. İnsanlar type, find, mklink gibi
komutları ezbere bilmiyorlar ama cat, grep, ln bir çok kişi tarafından
halihazırda benimsenmiş durumda. Ayrıca ssh, curl gibi şu anda komut satırında
yararlanılamayan araçlar da bir takım yazılımların dağıtılmasında anahtar görevi
görüyor.

Örnek:

```sh
$ ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
```

**Paket Yöneticisi:**\\
“Git” kurmak istiyorum. Ve bunu yalnızca komut satırından çalıştıracağım. Bunun
için msysgit kurmak durumunda kalıyorum. Oysa ki Debian’da olsam “apt-get
install git” yazdıktan sonra git kullanabilir hale gelecektim.

**Path Ayracı:**\\
Dosya pathleri için “/” ayracı artık Windows-dışı sistemlerde bir standart.
Windows bir noktaya kadar pathleri “/” ile kullandığımızda bir çevirim
gerçekleştiriyor. Fakat bunun Windows için de standart olması ne kaybettirir?
/media/c/Program Files/ benzeri bir yapı belki Linux’un avantajlarından biri
olan dizin-seviyesi partition desteği için ilk adım olabilir.

**CMD.EXE yerine Bash:**\\
Çok radikal bir fikir olmasa Windows da bash kullansın diyebilirdim. Geriye
dönük uyumluluk için CMD.EXE kalmaya devam eder, geliştirilmeye devam eden bash
ise bir port olarak Windows’un komut arayüzü olurdu. Böylece shell scriptler de
destekleniyor olurdu. Ama bunun yakın gelecekte olmasının mucize olduğunu
düşünüyorum.

Benim tespitlerim bu kadar.

Yazılımcılar olarak bir çok araç kullanıyor ve yazıyoruz. Bu araçların
kullanılması için komut satırı genelde önem arz ediyor. Modern web frameworkleri
komut satırından vazgeçmeyi bırakalım, tamamiyle komut satırı kullanımına
yönelmiş durumdalar. Bu araçları yazanlar Linux ve OS X için ortak komut
yazabiliyorken uyumsuzluk çıkaran Windows için bu komutları yazmaya tenezzül
dahi etmiyorlar.

Maddeleri “uygulanabilir”den “fantezi”ye doğru sıralamaya özen gösterdim. Bir
gün işin Windows tarafında da kendimi “eksik hissetmeyeceğim” günü iple
çekiyorum.
', 'tr          ', '01J6ZJ7QNBGTHZPRTAN02SPJXQ', '', 'Microsoft’a güçlü bir komut satırı için açık mektup');
insert into "story_tx" ("content", "locale_code", "story_id", "summary", "title") values ('Hello from news', 'en          ', '01J3YPPRB2EWZK7N31DMC50Y8C', 'It''s a greeting message from AYA''s brand new feature.', 'Hello world!');
insert into "story_tx" ("content", "locale_code", "story_id", "summary", "title") values ('Yazının öncülü olan
[DevRel (Developer Relations) Yolculuğum — Nasıl başladım?](./20221004-devrel-01)’ı
bu yazıdan önce okumanızı tavsiye ederim.

Ayrıca önceden iliştirmem gereken bir diğer not; bu yazı serisinde mümkün
olduğunca herhangi bir firmanın iç dinamiklerini açık etmeyeceğim. O nedenle bu
kısım en fazla herhangi bir iş mülakatında geçmiş tecrübe olarak duyabileceğiniz
kadar özet olacak.

## DevRel Nedir?

DevRel, şemsiyesinin altına birçok farklı disiplini alan çatı bir terim olduğu
için tek bir işlev tanımlamak yerine şemsiye altında gerçekleşen faaliyetlerle
tanımlamanın daha doğru olduğu bir disiplindir.

Ama illa bir tanım gerekiyorsa; “DevRel, bir organizasyonun ve geliştiricilerin
karşılıklı faydalarına hizmet etmesi için, strateji ve taktikler hazırlayan
organizasyon işlevidir” diyebiliriz.

Kavramsal ifadeler yeterince akıl karıştırdıysa, bu noktada bir es verip konuyu
örnekleyerek anlatmak isterim. “Halkla ilişkiler” nasıl “müşteri
memnuniyeti”nden tutalım “kriz yönetimine” kadar içine birçok işlevi
alabiliyorsa; “Geliştiricilerle ilişkiler” anlamına gelen “DevRel” de benzer bir
şekilde oldukça kapsayıcı bir terim.

Şunları tanımlamak gerekiyor:

- DevRel hangi kurum çatısı altında faaliyet gösteriyor?
- Bu kurum geliştiriciler ile nasıl bir ilişki içerisinde?
- Bu kurumun geliştiricilerle ortak paylaştığı alanlarda nasıl stratejileri ve
  hedefleri bulunuyor?

Dolayısıyla bugün AWS, Google, Microsoft v.b. firmaların DevRel’leri son
kullanıcıları geliştiriciler olan ürünlerin, platformların ve araçların
pazarlamalarında da etkin rol oynarken; bazı DevRel’lerin ne yaptığı konusunda
hiçbir fikriniz olmayabilir.

## Getir’de DevRel ne yapıyor? Ben ne yaptım?

Getir’de benim pozisyonum “DevRel Director” idi. Bu CTO’ya rapor eden bir
yönetici pozisyonu olarak tasarlandı. Daha önce emsali olmayan bir tanım ile
başladığımız için önce kitleyi ve hedefleri belirleyip konuştuk.

O günkü hedef geliştirici kitlesi içeride yer alması planlanan, global ölçekte
1000 kişiyi bulması düşünülen teknoloji departmanındaki ve teknoloji
departmanına katılmasını isteyebileceğimiz potansiyel insanlardı. Hızlı scale-up
etmiş, iletişim dili ingilizce olan, 4 farklı tech-hub’ı olan, global bir
teknoloji departmanında;

- Çalışan deneyimini arttırmak,
- İnsanların firmanın mühendislik kültürünü tanımasını ve bu kültürün bir
  parçası olması istemesini sağlamak,
- İyi mühendislik pratiklerinin benimsenmesini, konuşulmasını ve adaptasyonunu
  sağlamak,
- Çalışanların potansiyelini gerçekleştirebilmesini sağlayabilecek alanlar
  tanımlamak,
- Yöneticilere ve diğer departmanlara mühendislik kitlesi ile ilgili
  insight/içgörüler paylaşmak,

gibi stratejik hedefler bulunuyordu. Bu hedeflere de Technical Writer’dan
tutalım, Community Manager’a, hatta Open Source Developer’a uzanan çeşitli bir
şapkaların bulunduğu bir takım kurarak ilerlenilebilecek bir yol haritası
çıkartarak koşacaktık.

Kısacası DevRel Director tek tabanca “DevRel” olmak için değil, aslında firma
çatısı altında bir Developer Relations perspektifi kazandırmak ve stratejiler
belirlemek, bunu da birden fazla DevRel pozisyonuyla yapmak gibi konumlanmıştı.

10 aylık süre içerisinde sağlanan çalışmalardan açık olarak yazabileceklerimden
bazıları:

- 3 bootcamp tamamlandı, 14 kişi ekibe kazandırıldı
- Dış etkinliklerde konuşmalar yapıldı
- May 4th Star Wars günü kutlanmaya başlandı
- Getir Tech Talks ismiyle teknik konuşmaların olduğu bir etkinlik serisi
  başladı
- Hiring süreçleri ile ilgili iyileştirmeler üzerine çalışıldı
- Onboarding süreçleri ile ilgili iyileştirmeler üzerine çalışıldı
- Teknoloji departmanındaki roller için “gözlem-tabanlı” kariyer yolları
  oluşturuldu
- Uzun süre aynı görevde kalan kişiler için Internal Mobility Policy düzenlendi
- Teknoloji-Data-Product ekiplerinin squadlardan beklentilerini netleştirmek
  üzere “Maturity Matrix” çalışması ele alındı

Bu süreç içinde DevRel’liğin yanı sıra Agile Office’ın da sorumluluğunu alarak
ilerlediğim, çok şey öğrendiğim ve özveriyle çalıştığım keyifli bir süreç
geçirdim. Aynı süreç içerisinde firmanın decacorn olması gibi göğsümüzü kabartan
gelişmeler de yaşadık, pandemideki yasal kısıtlamalar nedeniyle ofise
gidemediğimiz de oldu.

## Bundan sonra nasıl ilerliyorum?

Yazının önceki bölümünde DevRel’in nereyse her kurumda olması gereken bir
“bakış” olduğunu savunduğumu, bu perspektifi de zaten birfiil mühendislik
yöneticisi kimliğim ile oluşturup daha sonra konsantre bir danışmanlığa
evirdiğimi anlatmıştım.

Ekim 2022 itibariyle artık Getir’den ayrılmış bulunuyorum. Ama bundan sonra
DevRel Director görevinde olmamam DevRel disiplinine yaptığım yatırımı yok
etmediği gibi; bu alanda halen bir şeyler üretmeme engel de değil.

Bu ay itibariyle direkt CTO ünvanıyla mühendislik yönetimine geri döndüğüm ancak
DevRel alanındaki biriktirdiğim know-how’ı yine işime yansıttığım bir düzende
ilerleyeceğim.

Aynı zamanda burada anlattıklarımı belki
[YouTube Yayınlarımda](https://eser.live) tekrar ele alabilirim.

Sağlıcakla kalın,
', 'tr          ', '01J6ZJ7SMXPB77HCZXSAQBN57W', 'Bu yazıda DevRel Director olarak Getir''de yaptıklarımı anlattım; mühendislik kültürünü güçlendirmek, çalışan deneyimini iyileştirmek ve teknik ekibe değer katmakla ilgili bakış açıları paylaşmaya çalıştım.', 'DevRel (Developer Relations) Yolculuğum — DevRel Director Olmak (2/2)');
insert into "story_tx" ("content", "locale_code", "story_id", "summary", "title") values ('Bugün Superpeer''in kapanmasının duyurusuyla birlikte Türkiye yazılım ekosisteminde aktif kullandığımız mecralardan birini kaybettiğimizi hissettim. Twitter''da [bu platformun açığını açık kaynak bir alternatifi ile destekleyebileceğimize ilişkin bir post](https://x.com/eserozvataf/status/1863810751246340317) paylaştım. Ve her zaman olduğu gibi buna karşılık da açıktan 1-2, DM''den onlarca mesaj aldım konuyla ilgili.

Ancak Batuhan''ın "cal.com da iş görüyor" fikrine ilişkin Superpeer''in "bence" yazılım ekosistemimizde hangi boşluğu doldurduğu ve neden cal.com ile bir tutulamayacağını düşündüğümü aktarmak, bu vesileyle eksiklere dem vurmaya başlamışken aklımdaki yazılımcı platformunu tariflemek istiyorum.

Öncelikle Superpeer''den başlayalım. Öncelikle yazılımcı platformu diyorum, ama superpeer daha çok benim için "yazılımcıların da kullandığı bir uzman platformu" idi. Evet burada o uzmanlara erişip ücretli/ücretsiz görüşme sağlayabiliyordunuz ancak katma değer kısmı bana kalırsa randevu ayarlama bölümü değil "uzman kişinin bir landing page''e sahip olması" kısmıydı. Farkındaysanız insanlar bu linkleri verirken bile "superpeer profilim" derken, diğer sistemlerde "cal.com/eserlive üzerinden randevu oluşturabilirsiniz" gibi lanse ediyorlar. Kısacası superpeer''i biz kartvizit ve landing page gibi kullanıyorduk.

Gelelim konunun mutlaka böyle bir sisteme neden ihtiyacımız var kısmına. Açık Yazılım Ağı sonrasında birçok kişiden "ne kadar çok işinin uzmanı, tanımadığımız veya yanlış tanıdığımız Türk Yazılımcı varmış" cümlesini duydum. Bu nedenle AYA''nın deprem faaliyetleri sonrasında ayağa kaldırmaya çalıştığım ilk web sitesi projesi AYA bünyesinde katkı sağlayan insanlara profil oluşturmaya imkan verecek bir web sitesi olmuştu. Ancak o dönem iş değişiklikleri gibi konulardan dolayı bunu canlıya alamadım.

Biz ekosistemdeki yazılımcıları maalesef twitter/x''deki click-bait gündemlerden, artık çoğu insanın çok temel şeyler konuşuluyor veya YouTube''da onlarca içerik var diye gitmeye üşendiği etkinliklerden tanıyabiliyoruz. Ama bunlar maalesef "yazılım" ve "üretim" odaklı değil tam tersi daha "lifestyle" temelli, paylaşımda yazılım bilgisinin değil yazılımcının odağında olduğu bir gündeme hapsediyor bu medium''lar bizi. Bu da sektörümüzü maalesef eskisine oranla meslektaşlarından pek hazetmeyen, bazen "yazılımcı olarak anılmak dahi istemediğimiz" bir noktaya getiriyor. Bir noktada bu mecraların vurguladığı noktaları değil, tam tersine "bizim bizzat kendimizi ifade edebildiğimiz" bir mecranın tüm bu zehirli, negatif ortamı nötrleştireceğine eminim.

Zaten kendini iyi ifade eden yazılımcılar ise genelde bir web sitesi oluşturmakla bunu sağlamaya çalışıyor. Şu an bu yazıyı okuduğunuz eser.live benim aynı zamanda içeriklerimi paylaştığım bir platform. Burası olmasa ne GitHub profilim ne Twitter''ım ne de YouTube''um benim hakkımda "genel derli toplu" bir bilgiyi size sunamayacak. Bence biz bu durumu o kadar içten fark ediyoruz ki bento, linktree gibi profiller oluşturup iz bırakmaya çalışıyoruz. Ama bana kalırsa halen bir şeyler eksik kalıyor. Ancak bu kişisel web siteleri ise bir süre sonra zamana yeniliyor, izlenebilirliği birer web sitesi olduklarından düşük kalıyor ve güncelleme almaz hale geliyorlar.

Hashnode, medium, substack ise çok yazı odaklılar. Biraz özelleştirmek, YouTube''da ürettiğiniz, attığınız tweet''i paylaşacağınız, GitHub''daki projenizi tanıtabileceğiniz hiçbir ensturman vermiyor size.

Bu nedenle ben eser.live gibi bir sistemin daha da öteye gidip (siteyi hazır buradayken biraz kurcalayabilirsiniz) "multi-tenant" bir uygulama olmasını, her yazılımcının kendini ve projelerini anlatabildiği landing page''ler oluşturabildiği, geçmişini paylaşabildiği, yazılarını, videolarını koyabildiği bir yer haline getirmek için acik.io domain''ini alıp (açık yazılım ağı''ndan hareketle) [hafif hafif bir şeyler kodlamaya başlamıştım](https://github.com/eser/acik.io/tree/template). Superpeer''in haberi ise şu an sürecimi hızlandırdı.

Eğer tasarım, Figma taraflarınız güçlüyse bence oturup bir hafta sonu workshop yapıp benzer bir şeyler çıkartalım. Hatta v0''da uzman olduğunu düşünen varsa direkt frontend''i çıkartalım, söz veriyorum backend tarafını ben hallederim :) Ulaşım için direkt tweet''in altına yazmanız yeterli.
', 'tr          ', '01JE63T76EBA8SPHPX9VJ9K1M0', 'Superpeer''in kapanmasıyla yazılım ekosistemimizde önemli bir boşluk oluştuğunu düşünüyorum. Yazılımcıların projelerini ve içeriklerini düzenli paylaşabileceği, kişisel sitelerin ve mevcut platformların eksiklerini kapatan, açık kaynak ve çok kullanıcılı yeni bir platform fikrini kaleme aldım.', 'Developer Platformu Eksikliği');
insert into "story_tx" ("content", "locale_code", "story_id", "summary", "title") values ('## Dinlemek, okumak ve bahsedileni anlamak

Birileri ile iletişim halindeyken onların ne anlatmaya çalıştığını anlamaya
çalışırım. Her metinin/konuşmanın bazı kusurları ve eksik noktaları vardır. O
eksik bırakılan yerlerden, söylemdeki hatalardan dolayı okumayı/dinlemeyi kesip
direkt yanıtlamaya çalışmam. Aktarılanı tekrar soğuk kanlılıkla değerlendirip
bana aktarılan fikirdeki/iddiadaki bakış açısını yakalamaya çalışırım. İçimdeki
muhakeme sonucu bahsedilen meseleyi benimsemek zorunda değilim ancak bunu bir
veri / yaklaşım olarak değerlendirmekte hiçbir sakınca olmaz.

## Saygı ve konuya odaklı kalmak

Ortada bir tartışma konusu olduğunda konuşmayı o kapsamın dışına taşırmamaya
odaklanırım. İletişimdeki kişilerin kim olduğu, ne yaptığı, ne düşündüğü, neye
inandığı/inanmadığı, neyi sevdiği/sevmediği, neyi istediği/istemediği, vb.
konuda sundukları argümanlarla ilintili değildir. Eğer biri ile tartışmaya
girdiysem, o kişiye saygı duymak zorunluluğum vardır. Aksi takdirde hiçbir
tarafın bir kazancının olamayacağı bir tartışmaya girmek anlamsızdır.

## Çeşitlilik ve kapsayıcılık

Toplumda dezavantajlı kesimler ve kişiler vardır. Bu noktada
insanları genellerken söylemlerden kapsam dışında bıracak,
onları dışlanmış hissettirecek, gücendirecek ifadelerden kaçınırım.
', 'tr          ', '01J6ZJ7P393E92TMR12K3GY0FS', 'İletişim kurarken benimsediğim ilkeler üzerine ilkelerimi aktarmaya çalıştım.', 'İletişim İlkelerim');
insert into "story_tx" ("content", "locale_code", "story_id", "summary", "title") values (E'(Yazı halen geliştiriliyor, şundan da bahset derseniz
[https://eser.live/qa](https://eser.live/qa) üzerinden bana ulaşabilirsiniz)

## Motivasyon

Superpeer ve benzeri kanallardan benimle iletişime geçen birçok insanla benzer
temalarda konuşmam üzerine bu yazıyı yazma kararı aldım.

Normal şartlarda yüksek lisansını Eğitim Fakültesi’nde tamamlamış, yine
lisansüstü eğitimimde eğitim ve psikoloji ile ilgili birçok akademik bilgiye
sahip olmuş biri olarak, birçok insanın “ne bilmediğini bilmediği”ni düşündüğüm
alanlar “öğrenim” ve “gelişim” konuları. Bu gelişim ister kariyer, ister kişisel
ister başka temada olsun fark etmiyor. Öğrenim psikolojik ve kişiye özgü bir
kavram. Tek bir doğru olmadığı gibi sürece deterministik bakmak birçok
zararı beraberinde getiriyor. Bu nedenle beni genellikle “başarmak için bu
adımları takip edin” derken pek duymazsınız, daha empirist/deneyci yaklaşımlarım
vardır.

Peki ben tavsiye vermekten kaçıyorsam ne için bu yazıyı yazıyorum? Harika bir
kaos teorisi tutkunu olduğum için elbette! Bu yazı sayesinde hem kaoscu ve
bilinmezci yaklaşımlarımla hava atacağım hem de gözlemlediğim kadarıyla öğrenme
yoluna girmiş insanların kendini yetersiz hissettiklerinde bir motivasyon
bulmaları için kendi çapımda bir katkı sunmaya çalışacağım. Elbette ne kadar
başarılı olurum bilinmez, bir iddiam bulunmuyor.

# Yol Haritası Bunun Neresinde?

Evet, bu halen bir yol haritası. Burada yapmaya çalışacağım şey şu olacak:
Birçok insanın kendine koyduğunu gözlemlediğim (naçizane), sık rastlanan mental
engelleri tanımlamaya çalışacağım.

Dolayısıyla siz kendi yolunuzu kendiniz bulurken (bence olması gereken bu) size
yolda rastlayabileceğiniz, sizi yolunuzdan döndürecek veya yolunuzu
uzatacak/zorlaştıracak tehlikelerden bahsedeceğim.

Bunun için istemeden de olsa tavsiye verir tonda konuşacağım ama bunları benim
kendi kendi kendime problem çözme yöntemlerim gibi ele alırsanız sevinirim.

## Yolda Karşılaştığımız Tehlikeler

### Genellemelere Maruz Kalmak

Hepimiz hayatımızın farklı alanlarında da olsa insanların genellemelerine ve
tespitlerine olumsuz yönde maruz kalmışızdır.

Örneğin 6 ayda react öğrenebilirsiniz denilmiştir, ama siz öğrenememişsinizdir:

- bu tespitlerin bir ortalama, varsayımlar ve korelasyon ile oluştuğunu,

- genelleme sahiplerinin sandığınız kadar ince eleyip sık dokumadan bu
  tespitleri 5-10 dakikalık düşünsel süreçlerle ortaya koyduğunu,

bilmekte yarar var.

Bu durumlarla karşılaştığımızda genelleme/iddia sahibinin ne anlatmaya
çalıştığına, -eğer burada bizim için geçerli bir durum söz konusu ise- buradan
nasıl bir çıkarım/fikir elde edebileceğimize odaklanmalıyız.

### Reçetelere Bağlı Kalmak

Bir plan olması her zaman iyiye işaret olsa da planlar değişebilir ve
kaybolmamak için bize yardımcı olan unsurlardır. Planları amaç olarak değil araç
olarak konumlandırmak dolayısıyla planlar gerçekleşmediğinde başarısız olmuş
hissiyatıyla toptan bir vazgeçişe yönlenmek kendimizi içinde bulabileceğimiz bir
tuzak.

Dolayısıyla roadmap.sh gibi sitelerde size verilen reçetelerdeki
kavramları bilmeniz elbette size zarar vermez, ancak bu sitelerin reçetelerini
harfiyen uygulamak yerine bir şeyler yaparken ihtiyacınızın olduğunu fark
ettiğiniz kavramlara yoğunlaşmayı deneyebilirsiniz. HTTP nedir’i ve protokolün
detaylarını sonradan öğrendiğinizde dünyanın sonu gelmiyor.

### Süreci Yok Saymak

Kendinizi X yıldır sektörde bulunan birisiyle kıyaslayıp, “Ben kesinlikle
bunları öğrenemem” diye vazgeçmek ve diğer insanların sanki her şeyi bir süreç
içinde öğrenmek yerine bir gece aniden onlara vahiy geldiğini düşünmek de
sıklıkla karşılaşılan durumlar arasında.

Ancak bir az üzerine düşünürsek, bu insani olarak çok acımasızca. İnsan olarak
İlköğretimden itibaren on yıllarca gördüğümüz matematik ve ana dilimizle ilgili
eğitim sürecimizin sonucunda geometri, dil bilgisi gibi kavramları temel
düzeylerde konuşabilecek ve problem çözebilecek bir noktaya erişebiliyoruz.

Yazılım da öyle, uzun bir sürecin ertesinde eğitimimiz öğrenim ve pratiğe
dönüşüyor.

### Başarı Listesini Tutmamak

Öğrenme bir sürece yayıldığı için kimi zaman biz ısıya maruz kalan kurbağa gibi
durumumuza “alışıyoruz”, 2 ay önce bilmediğimiz birçok şeyi bilsek dahi, henüz
“nihai hedefe” ulaşamamış olmak bizi halen “başarısız” hissettiriyor.

Eğer çok güçlü duygusal sermayemiz, metanetli bir karakterimiz yoksa vazgeçmeye
çok açık olabiliyoruz.

### Yapabildiklerini ve Öğrendiklerini Küçümsemek

Öncelikle geçmiş olsun, deneyimli bir profesyonel olsanız dahi ileride de bu
hissiyatınızdan asla arınamayabilirsiniz. Hatta bunu ara ara duyacağınız
“Imposter Syndrome” ile ile ilişkilendirebiliriz belki.

Öğrendiğiniz birçok kavram emin olun boşuna değildir. Ancak siz bu öğrenimleri
bir pratik projede kullanmayıp, diğer insanların bu kadar temel bilgilere bile
ara sıra hakim olmadığını veya bunları kullanamadığını deneyimlemediğiniz için
size “fasa fiso” geliyor olabilir.

Örneğin exercism üzerinde linked list problemlerini kolayca çözen birileri bizim
senior software engineer’ların live interview’ında da bu tarz sorular sorduğumuz
bilgisine sahip olmadığı için kendi bilgisini çok basite indirgeyebiliyor.

### Diğer İnsanların Acı Çekmediğini Düşünmek

Sizin tanık olmadığınız tüm “sonuç”lar, dışarıdan her şeyi mükemmel ve kusursuz
görmenize sebep olabiliyor. Herkes diş bakımı için dişçi koltuklarında nice
acılar çeker, her gün sabah akşam dişlerini temiz tutmak için uğraşır hatta çoğu
insan cerrahi operasyonlara maruz kalır ancak biz Instagram’da onları görüp “Aa
ne kadar güzel dişleri var, tanrı vergisi” deriz.

Bu bir yanılsama ve dışarıdan “tamamlanmış” görünen birçok hikaye
geliştirilirken birçok başarısızlıkla karşılaşılıyor. Kimse bir anda bir noktaya
gelmiyor, kimse yazdığı projeyi anında çalıştıramıyor. Geliştirme sürecince
birçok hata alıyor ve ufak bir stil düzenlemesi için günler kaybedebiliyor.

### Diğer İnsanlarla Karşılaştırmak

\\-

### Diğer İnsanları Fazla Dinlemek

\\-
', 'tr          ', '01J6ZJ7PZJWHBFYY8Y5A0FSM49', 'Birçok kişinin öğrenme sürecinde karşılaştığı motivasyon sorunlarına değinmek istedim. Yol haritalarına katı şekilde bağlı kalmak, genellemelere takılmak, süreci göz ardı etmek ve başarıları küçümsemek gibi mental engellerin farkında olmak, gelişim yolunda karşılaşılan kritik engeller.', 'Öğrenme Yol Haritası (Bilhassa Junior ve Yazılımcılar İçin)');
insert into "story_tx" ("content", "locale_code", "story_id", "summary", "title") values (E'Bu yazıda bana birçok kez sorulan bilim kurgu janrası ile tanışmak için
öncelikle neler izlemek gerekli sorusuna yanıt olarak verdiğim, kendi kişisel
öneri listemi paylaşacağım.

Liste oluşturulurken dikkat ettiğim unsur, önerilen eserlerin mükemmeliyetinden
çok, işledikleri konu/temalar nedeniyle bilim kurgu eserlerine giriş mahiyeti
oluşturması ve janrı tanıtmasına odaklı. Eğer önerileri takip edecekseniz
şimdiden keyifli seyirler :)

Bilinmesi gereken kavramlar:

- Cyberpunk, Siberpunk, Siber çılgınlık\\
  [https://tr.wikipedia.org/w/index.php?title=Siberpunk](https://tr.wikipedia.org/w/index.php?title=Siberpunk&useskin=vector)

- Distopya\\
  [https://tr.wikipedia.org/w/index.php?title=Distopya](https://tr.wikipedia.org/w/index.php?title=Distopya)

Filmler:

- Blade Runner

- The Matrix

- Dune

- Moon

- Alien

- Terminator

- The Man from Earth

- Back to the Future

- V for Vendetta

- A Clockwork Orange

- 1984

- Upgrade

- Snowpiercer

Diziler:

- Star Trek

- X-Files

- Fringe

- Battlestar Galactica

- Twilight Zone

- Person of Interest

- Doctor Who

- Serial Experiments Lain

- Westworld
', 'tr          ', '01J6ZJ7S2487FJ2Q1MZB95ZYHE', 'Bir bilim kurgu hayranı ve takipçisi olarak bilim kurgu dünyasına ilgi duyanlar ve yeni başlayanlar için bir liste oluşturmaya çalıştım.', 'Bilim Kurgu / Distopya Eserlerle Tanışma Öneri Listesi');

insert into "story_publication" ("created_at", "deleted_at", "id", "kind", "profile_id", "properties", "story_id", "updated_at") values ('2024-09-04 22:24:46.979119+00', NULL, '01J6ZJ7SSNZ0JT71YTCB1YJ8F2', '', '01J27JA3WM81T7D1TCNJ6Q06P5', NULL, '01J6ZJ7SSNZ0JT71YTCB1YJ8F2', NULL);
insert into "story_publication" ("created_at", "deleted_at", "id", "kind", "profile_id", "properties", "story_id", "updated_at") values ('2024-09-04 22:24:44.520546+00', NULL, '01J6ZJ7QCTXKKS7DJ9NYMC9ZSH', '', '01J27JA3WM81T7D1TCNJ6Q06P5', NULL, '01J6ZJ7QCTXKKS7DJ9NYMC9ZSH', NULL);
insert into "story_publication" ("created_at", "deleted_at", "id", "kind", "profile_id", "properties", "story_id", "updated_at") values ('2024-09-04 22:24:46.827225+00', NULL, '01J6ZJ7SMXPB77HCZXSAQBN57W', '', '01J27JA3WM81T7D1TCNJ6Q06P5', NULL, '01J6ZJ7SMXPB77HCZXSAQBN57W', NULL);
insert into "story_publication" ("created_at", "deleted_at", "id", "kind", "profile_id", "properties", "story_id", "updated_at") values ('2024-09-04 22:24:42.741007+00', NULL, '01J6ZJ7NN6WHYB587G2TXRTS3Q', '', '01J27JA3WM81T7D1TCNJ6Q06P5', NULL, '01J6ZJ7NN6WHYB587G2TXRTS3Q', NULL);
insert into "story_publication" ("created_at", "deleted_at", "id", "kind", "profile_id", "properties", "story_id", "updated_at") values ('2024-09-04 22:24:41.899933+00', NULL, '01J6ZJ7MTZNZ207P1NYAA41A61', '', '01J27JA3WM81T7D1TCNJ6Q06P5', NULL, '01J6ZJ7MTZNZ207P1NYAA41A61', NULL);
insert into "story_publication" ("created_at", "deleted_at", "id", "kind", "profile_id", "properties", "story_id", "updated_at") values ('2024-09-04 22:24:42.877668+00', NULL, '01J6ZJ7NSFX726JYRV59WX72MP', '', '01J27JA3WM81T7D1TCNJ6Q06P5', NULL, '01J6ZJ7NSFX726JYRV59WX72MP', NULL);
insert into "story_publication" ("created_at", "deleted_at", "id", "kind", "profile_id", "properties", "story_id", "updated_at") values ('2024-09-04 22:24:42.461191+00', NULL, '01J6ZJ7NCGNHWTP081Q3DD0BXN', '', '01J27JA3WM81T7D1TCNJ6Q06P5', NULL, '01J6ZJ7NCGNHWTP081Q3DD0BXN', NULL);
insert into "story_publication" ("created_at", "deleted_at", "id", "kind", "profile_id", "properties", "story_id", "updated_at") values ('2024-09-04 22:24:42.037763+00', NULL, '01J6ZJ7MZ7DF5QEDZQFM7K8NV5', '', '01J27JA3WM81T7D1TCNJ6Q06P5', NULL, '01J6ZJ7MZ7DF5QEDZQFM7K8NV5', NULL);
insert into "story_publication" ("created_at", "deleted_at", "id", "kind", "profile_id", "properties", "story_id", "updated_at") values ('2024-09-04 22:24:43.37063+00', NULL, '01J6ZJ7P8W06V3WY2R70S53E8W', '', '01J27JA3WM81T7D1TCNJ6Q06P5', NULL, '01J6ZJ7P8W06V3WY2R70S53E8W', NULL);
insert into "story_publication" ("created_at", "deleted_at", "id", "kind", "profile_id", "properties", "story_id", "updated_at") values ('2024-09-04 22:24:46.537267+00', NULL, '01J6ZJ7SBVZ2NHN2BE00GT00ZF', '', '01J27JA3WM81T7D1TCNJ6Q06P5', NULL, '01J6ZJ7SBVZ2NHN2BE00GT00ZF', NULL);
insert into "story_publication" ("created_at", "deleted_at", "id", "kind", "profile_id", "properties", "story_id", "updated_at") values ('2024-09-04 22:24:45.513668+00', NULL, '01J6ZJ7RBVQ3GD7H1NHN0ZYBSY', '', '01J27JA3WM81T7D1TCNJ6Q06P5', NULL, '01J6ZJ7RBVQ3GD7H1NHN0ZYBSY', NULL);
insert into "story_publication" ("created_at", "deleted_at", "id", "kind", "profile_id", "properties", "story_id", "updated_at") values ('2024-09-04 22:24:44.096282+00', NULL, '01J6ZJ7PZJWHBFYY8Y5A0FSM49', '', '01J27JA3WM81T7D1TCNJ6Q06P5', NULL, '01J6ZJ7PZJWHBFYY8Y5A0FSM49', NULL);
insert into "story_publication" ("created_at", "deleted_at", "id", "kind", "profile_id", "properties", "story_id", "updated_at") values ('2024-09-04 22:24:43.686257+00', NULL, '01J6ZJ7PJR08D72KM5BAF67808', '', '01J27JA3WM81T7D1TCNJ6Q06P5', NULL, '01J6ZJ7PJR08D72KM5BAF67808', NULL);
insert into "story_publication" ("created_at", "deleted_at", "id", "kind", "profile_id", "properties", "story_id", "updated_at") values ('2024-09-04 22:24:44.232534+00', NULL, '01J6ZJ7Q3VH05VHF4889WWE29V', '', '01J27JA3WM81T7D1TCNJ6Q06P5', NULL, '01J6ZJ7Q3VH05VHF4889WWE29V', NULL);
insert into "story_publication" ("created_at", "deleted_at", "id", "kind", "profile_id", "properties", "story_id", "updated_at") values ('2024-09-04 22:24:45.10021+00', NULL, '01J6ZJ7QYWGM91G0Q6HQRPPPES', '', '01J27JA3WM81T7D1TCNJ6Q06P5', NULL, '01J6ZJ7QYWGM91G0Q6HQRPPPES', NULL);
insert into "story_publication" ("created_at", "deleted_at", "id", "kind", "profile_id", "properties", "story_id", "updated_at") values ('2024-09-04 22:24:43.956782+00', NULL, '01J6ZJ7PV71HMMS0R23NCRMKGT', '', '01J27JA3WM81T7D1TCNJ6Q06P5', NULL, '01J6ZJ7PV71HMMS0R23NCRMKGT', NULL);
insert into "story_publication" ("created_at", "deleted_at", "id", "kind", "profile_id", "properties", "story_id", "updated_at") values ('2024-09-04 22:24:43.510896+00', NULL, '01J6ZJ7PD8J8YP66QE91Q2DH5Q', '', '01J27JA3WM81T7D1TCNJ6Q06P5', NULL, '01J6ZJ7PD8J8YP66QE91Q2DH5Q', NULL);
insert into "story_publication" ("created_at", "deleted_at", "id", "kind", "profile_id", "properties", "story_id", "updated_at") values ('2024-09-04 22:24:44.924375+00', NULL, '01J6ZJ7QSGP0K3EDZM0CNCKN74', '', '01J27JA3WM81T7D1TCNJ6Q06P5', NULL, '01J6ZJ7QSGP0K3EDZM0CNCKN74', NULL);
insert into "story_publication" ("created_at", "deleted_at", "id", "kind", "profile_id", "properties", "story_id", "updated_at") values ('2024-09-04 22:24:42.313649+00', NULL, '01J6ZJ7N7VSV25K6ZD7ZVMQ0KV', '', '01J27JA3WM81T7D1TCNJ6Q06P5', NULL, '01J6ZJ7N7VSV25K6ZD7ZVMQ0KV', NULL);
insert into "story_publication" ("created_at", "deleted_at", "id", "kind", "profile_id", "properties", "story_id", "updated_at") values ('2024-09-04 22:24:45.373847+00', NULL, '01J6ZJ7R7GBWBFNK9ME20DJEVE', '', '01J27JA3WM81T7D1TCNJ6Q06P5', NULL, '01J6ZJ7R7GBWBFNK9ME20DJEVE', NULL);
insert into "story_publication" ("created_at", "deleted_at", "id", "kind", "profile_id", "properties", "story_id", "updated_at") values ('2024-09-04 22:24:44.654873+00', NULL, '01J6ZJ7QH103QMEQZK0QAQWJWY', '', '01J27JA3WM81T7D1TCNJ6Q06P5', NULL, '01J6ZJ7QH103QMEQZK0QAQWJWY', NULL);
insert into "story_publication" ("created_at", "deleted_at", "id", "kind", "profile_id", "properties", "story_id", "updated_at") values ('2024-09-04 22:24:42.593884+00', NULL, '01J6ZJ7NGMK53TAMW0C783FV0B', '', '01J27JA3WM81T7D1TCNJ6Q06P5', NULL, '01J6ZJ7NGMK53TAMW0C783FV0B', NULL);
insert into "story_publication" ("created_at", "deleted_at", "id", "kind", "profile_id", "properties", "story_id", "updated_at") values ('2024-09-04 22:24:45.653402+00', NULL, '01J6ZJ7RG7RCMAFHWRX01BCDEM', '', '01J27JA3WM81T7D1TCNJ6Q06P5', NULL, '01J6ZJ7RG7RCMAFHWRX01BCDEM', NULL);
insert into "story_publication" ("created_at", "deleted_at", "id", "kind", "profile_id", "properties", "story_id", "updated_at") values ('2024-09-04 22:24:41.627548+00', NULL, '01J6ZJ7MJDAG3GZHY49SE509RD', '', '01J27JA3WM81T7D1TCNJ6Q06P5', NULL, '01J6ZJ7MJDAG3GZHY49SE509RD', NULL);
insert into "story_publication" ("created_at", "deleted_at", "id", "kind", "profile_id", "properties", "story_id", "updated_at") values ('2024-09-04 22:24:43.819129+00', NULL, '01J6ZJ7PPXDJMG7H8M8K2C1BTF', '', '01J27JA3WM81T7D1TCNJ6Q06P5', NULL, '01J6ZJ7PPXDJMG7H8M8K2C1BTF', NULL);
insert into "story_publication" ("created_at", "deleted_at", "id", "kind", "profile_id", "properties", "story_id", "updated_at") values ('2024-09-04 22:24:46.687522+00', NULL, '01J6ZJ7SGH8HBPKXMWD3G0JGW9', '', '01J27JA3WM81T7D1TCNJ6Q06P5', NULL, '01J6ZJ7SGH8HBPKXMWD3G0JGW9', NULL);
insert into "story_publication" ("created_at", "deleted_at", "id", "kind", "profile_id", "properties", "story_id", "updated_at") values ('2024-09-04 22:24:46.225518+00', NULL, '01J6ZJ7S2487FJ2Q1MZB95ZYHE', '', '01J27JA3WM81T7D1TCNJ6Q06P5', NULL, '01J6ZJ7S2487FJ2Q1MZB95ZYHE', NULL);
insert into "story_publication" ("created_at", "deleted_at", "id", "kind", "profile_id", "properties", "story_id", "updated_at") values ('2024-09-04 22:24:44.385626+00', NULL, '01J6ZJ7Q8MYTANHEV0THKJ2XVN', '', '01J27JA3WM81T7D1TCNJ6Q06P5', NULL, '01J6ZJ7Q8MYTANHEV0THKJ2XVN', NULL);
insert into "story_publication" ("created_at", "deleted_at", "id", "kind", "profile_id", "properties", "story_id", "updated_at") values ('2024-09-04 22:24:47.112296+00', NULL, '01J6ZJ7SXTDJA9SZ9WGV7RD03H', '', '01J27JA3WM81T7D1TCNJ6Q06P5', NULL, '01J6ZJ7SXTDJA9SZ9WGV7RD03H', NULL);
insert into "story_publication" ("created_at", "deleted_at", "id", "kind", "profile_id", "properties", "story_id", "updated_at") values ('2024-09-04 22:24:45.235234+00', NULL, '01J6ZJ7R34NJXYD7X854XCH5AG', '', '01J27JA3WM81T7D1TCNJ6Q06P5', NULL, '01J6ZJ7R34NJXYD7X854XCH5AG', NULL);
insert into "story_publication" ("created_at", "deleted_at", "id", "kind", "profile_id", "properties", "story_id", "updated_at") values ('2024-09-04 22:24:46.087131+00', NULL, '01J6ZJ7RXSVEDNK55KKZZ8VZ2N', '', '01J27JA3WM81T7D1TCNJ6Q06P5', NULL, '01J6ZJ7RXSVEDNK55KKZZ8VZ2N', NULL);
insert into "story_publication" ("created_at", "deleted_at", "id", "kind", "profile_id", "properties", "story_id", "updated_at") values ('2024-09-04 22:24:41.760573+00', NULL, '01J6ZJ7MPJHY28R4M7KWF28SJE', '', '01J27JA3WM81T7D1TCNJ6Q06P5', NULL, '01J6ZJ7MPJHY28R4M7KWF28SJE', NULL);
insert into "story_publication" ("created_at", "deleted_at", "id", "kind", "profile_id", "properties", "story_id", "updated_at") values ('2024-09-04 22:24:45.79335+00', NULL, '01J6ZJ7RMK5Z7ETSNTHJYAVZ3Z', '', '01J27JA3WM81T7D1TCNJ6Q06P5', NULL, '01J6ZJ7RMK5Z7ETSNTHJYAVZ3Z', NULL);
insert into "story_publication" ("created_at", "deleted_at", "id", "kind", "profile_id", "properties", "story_id", "updated_at") values ('2024-09-04 22:24:42.174582+00', NULL, '01J6ZJ7N3G5HWRNDZ8YRJD8CDJ', '', '01J27JA3WM81T7D1TCNJ6Q06P5', NULL, '01J6ZJ7N3G5HWRNDZ8YRJD8CDJ', NULL);
insert into "story_publication" ("created_at", "deleted_at", "id", "kind", "profile_id", "properties", "story_id", "updated_at") values ('2024-09-04 22:24:43.014109+00', NULL, '01J6ZJ7NXRZC49DCQFFK8RV16R', '', '01J27JA3WM81T7D1TCNJ6Q06P5', NULL, '01J6ZJ7NXRZC49DCQFFK8RV16R', NULL);
insert into "story_publication" ("created_at", "deleted_at", "id", "kind", "profile_id", "properties", "story_id", "updated_at") values ('2024-09-04 22:24:41.469413+00', NULL, '01J6ZJ7KSEQ6ERWW48T3BZCQRY', '', '01J27JA3WM81T7D1TCNJ6Q06P5', NULL, '01J6ZJ7KSEQ6ERWW48T3BZCQRY', NULL);
insert into "story_publication" ("created_at", "deleted_at", "id", "kind", "profile_id", "properties", "story_id", "updated_at") values ('2024-09-04 22:24:46.4024+00', NULL, '01J6ZJ7S7MEQH7T95RG9NBBYFY', '', '01J27JA3WM81T7D1TCNJ6Q06P5', NULL, '01J6ZJ7S7MEQH7T95RG9NBBYFY', NULL);
insert into "story_publication" ("created_at", "deleted_at", "id", "kind", "profile_id", "properties", "story_id", "updated_at") values ('2024-09-04 22:24:45.950561+00', NULL, '01J6ZJ7RSFJV6VBYPT467BKFQ0', '', '01J27JA3WM81T7D1TCNJ6Q06P5', NULL, '01J6ZJ7RSFJV6VBYPT467BKFQ0', NULL);
insert into "story_publication" ("created_at", "deleted_at", "id", "kind", "profile_id", "properties", "story_id", "updated_at") values ('2024-09-04 22:24:43.191998+00', NULL, '01J6ZJ7P393E92TMR12K3GY0FS', '', '01J27JA3WM81T7D1TCNJ6Q06P5', NULL, '01J6ZJ7P393E92TMR12K3GY0FS', NULL);
insert into "story_publication" ("created_at", "deleted_at", "id", "kind", "profile_id", "properties", "story_id", "updated_at") values ('2024-09-04 22:24:44.792592+00', NULL, '01J6ZJ7QNBGTHZPRTAN02SPJXQ', '', '01J27JA3WM81T7D1TCNJ6Q06P5', NULL, '01J6ZJ7QNBGTHZPRTAN02SPJXQ', NULL);
insert into "story_publication" ("created_at", "deleted_at", "id", "kind", "profile_id", "properties", "story_id", "updated_at") values ('2024-12-03 14:18:21.353+00', NULL, '01JE63T76EBA8SPHPX9VJ9K1M0', '', '01J27JA3WM81T7D1TCNJ6Q06P5', NULL, '01JE63T76EBA8SPHPX9VJ9K1M0', NULL);
insert into "story_publication" ("created_at", "deleted_at", "id", "kind", "profile_id", "properties", "story_id", "updated_at") values ('2024-07-29 07:36:36.768032+00', NULL, '01J3YPPRB2EWZK7N31DMC50Y8C', '', '01JT4FS0K93K4GG4Z40ZDFY3EQ', NULL, '01J3YPPRB2EWZK7N31DMC50Y8C', NULL);

-- users

insert into "user" ("bsky_handle", "bsky_remote_id", "created_at", "deleted_at", "email", "github_handle", "github_remote_id", "id", "individual_profile_id", "kind", "name", "phone", "updated_at", "x_handle", "x_remote_id") values (NULL, NULL, '2025-07-13 17:23:16.244933+00', NULL, 'eser@acikyazilim.com', 'eser', '866558', '01K02CZZ4M5QBC350BJMBXBJFY', '01J27JA3WM81T7D1TCNJ6Q06P5', 'regular', 'Eser Ozvataf', NULL, NULL, NULL, NULL);
