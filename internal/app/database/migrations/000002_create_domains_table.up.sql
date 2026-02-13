CREATE TABLE domains (
    id BIGSERIAL PRIMARY KEY,
    batch_id INT NOT NULL,
    full_url VARCHAR(50) NOT NULL,
    status_code VARCHAR(50),
    redirect_url VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    result VARCHAR(50),

    CONSTRAINT fk_domain_domain_check_history 
        FOREIGN KEY (batch_id) REFERENCES domain_check_histories(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

CREATE INDEX idx_domain_batch_id ON domains(batch_id);