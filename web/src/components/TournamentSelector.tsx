import React, { useState } from 'react';
import './TournamentSelector.css';

interface TournamentSelectorProps {
  value: string;
  onChange: (tournament: string) => void;
  onSubmit?: () => void;
}

export function TournamentSelector({ value, onChange, onSubmit }: TournamentSelectorProps) {
  const [input, setInput] = useState(value);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onChange(input);
    onSubmit?.();
  };

  return (
    <form className="tournament-selector" onSubmit={handleSubmit}>
      <label htmlFor="tournament">Tournament</label>
      <div className="tournament-input-group">
        <input
          id="tournament"
          type="text"
          value={input}
          onChange={e => setInput(e.target.value)}
          placeholder="octagon"
        />
        <button type="submit">Load</button>
      </div>
    </form>
  );
}
