CREATE TABLE gerund (
    infinitive character varying NOT NULL PRIMARY KEY,
    gerund character varying NOT NULL,
    gerund_english character varying
);
CREATE TABLE infinitive (
    infinitive character varying NOT NULL PRIMARY KEY,
    infinitive_english character varying
);
CREATE TABLE mood (
    mood character varying NOT NULL PRIMARY KEY,
    mood_english character varying
);
CREATE TABLE pastparticiple (
    infinitive character varying NOT NULL PRIMARY KEY,
    pastparticiple character varying NOT NULL,
    pastparticiple_english character varying
);
CREATE TABLE tense (
    tense character varying NOT NULL PRIMARY KEY,
    tense_english character varying
);
CREATE TABLE verbs (
    infinitive character varying NOT NULL,
    mood character varying NOT NULL,
    tense character varying NOT NULL,
    verb_english character varying,
    form_1s character varying,
    form_2s character varying,
    form_3s character varying,
    form_1p character varying,
    form_2p character varying,
    form_3p character varying,
    PRIMARY KEY (infinitive, mood, tense)
);
