CREATE TABLE domains (
    id BIGSERIAL PRIMARY KEY,
    batch_id INT NOT NULL,
    full_url TEXT NOT NULL,
    status_code INT,
    redirect_url TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_domain_domain_check_history 
        FOREIGN KEY (batch_id) REFERENCES domain_check_histories(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

CREATE INDEX idx_domain_batch_id ON domains(batch_id);