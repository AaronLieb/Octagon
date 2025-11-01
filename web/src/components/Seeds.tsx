import React, { useState } from 'react';

interface SeedResult {
  name: string;
  rating: number;
  seed: number;
  id: number;
}

interface SeedsResponse {
  tournament: string;
  event: string;
  seeds: SeedResult[];
}

const Seeds: React.FC = () => {
  const [data, setData] = useState<SeedsResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [tournament, setTournament] = useState('octagon');
  const [redemption, setRedemption] = useState(false);
  const [showConfirmation, setShowConfirmation] = useState(false);
  const [publishing, setPublishing] = useState(false);

  const runSeeding = async (tournamentSlug: string, isRedemption: boolean) => {
    setLoading(true);
    setError(null);
    
    try {
      const response = await fetch('http://localhost:8080/api/seed', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          tournament: tournamentSlug,
          redemption: isRedemption,
        }),
      });
      
      if (!response.ok) {
        throw new Error('Failed to run seeding');
      }
      const result = await response.json();
      setData(result);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred');
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    runSeeding(tournament, redemption);
  };

  const publishSeeds = async () => {
    if (!data) return;
    
    setPublishing(true);
    setError(null);
    
    try {
      const response = await fetch('http://localhost:8080/api/seed/publish', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          tournament: tournament,
          redemption: redemption,
          players: data.seeds,
        }),
      });
      
      if (!response.ok) {
        throw new Error('Failed to publish seeding');
      }
      
      setShowConfirmation(false);
      alert('Seeding published successfully!');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred');
    } finally {
      setPublishing(false);
    }
  };

  const confirmSeeding = () => {
    setShowConfirmation(false);
    runSeeding(tournament, redemption);
  };

  const cancelSeeding = () => {
    setShowConfirmation(false);
  };

  if (loading) {
    return (
      <div className="loading">
        Running seeding algorithm...
      </div>
    );
  }

  return (
    <div className="container">
      <h1 className="title">Generate Tournament Seeding</h1>
      
      <form onSubmit={handleSubmit} className="form">
        <input
          type="text"
          value={tournament}
          onChange={(e) => setTournament(e.target.value)}
          placeholder="Tournament slug (e.g., octagon)"
          className="input"
        />
        <label style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
          <input
            type="checkbox"
            checked={redemption}
            onChange={(e) => setRedemption(e.target.checked)}
          />
          Redemption Bracket
        </label>
        <button type="submit" className="button">
          Generate Seeding
        </button>
      </form>

      {showConfirmation && (
        <div style={{
          position: 'fixed',
          top: 0,
          left: 0,
          right: 0,
          bottom: 0,
          backgroundColor: 'rgba(0,0,0,0.5)',
          display: 'flex',
          justifyContent: 'center',
          alignItems: 'center',
          zIndex: 1000
        }}>
          <div style={{
            backgroundColor: 'white',
            padding: '24px',
            borderRadius: '8px',
            maxWidth: '400px',
            textAlign: 'center'
          }}>
            <h3>Publish Seeding to start.gg?</h3>
            <p>This will update the seeding on start.gg for <strong>{data?.tournament}</strong>.</p>
            {redemption && <p><em>This will update the redemption bracket seeding.</em></p>}
            <div style={{ marginTop: '16px', display: 'flex', gap: '8px', justifyContent: 'center' }}>
              <button onClick={publishSeeds} className="button" disabled={publishing}>
                {publishing ? 'Publishing...' : 'Yes, Publish'}
              </button>
              <button onClick={() => setShowConfirmation(false)} className="button-secondary" disabled={publishing}>
                Cancel
              </button>
            </div>
          </div>
        </div>
      )}

      {error && (
        <div className="error">
          Error: {error}
        </div>
      )}

      {data && (
        <>
          <div className="tournament-info">
            <h2 className="tournament-name">{data.tournament}</h2>
            <p className="attendee-count">{data.seeds.length} players seeded for {data.event}</p>
          </div>

          <div className="table-container">
            <table className="table">
              <thead>
                <tr>
                  <th>Seed</th>
                  <th>Gamer Tag</th>
                  <th>Rating</th>
                </tr>
              </thead>
              <tbody>
                {data.seeds.map((seed) => (
                  <tr key={seed.seed}>
                    <td className="gamer-tag">{seed.seed}</td>
                    <td className="name">{seed.name}</td>
                    <td className="player-id">{seed.rating.toFixed(2)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          <div style={{ marginTop: '16px', textAlign: 'center' }}>
            <button onClick={() => setShowConfirmation(true)} className="button">
              Publish to start.gg
            </button>
          </div>
        </>
      )}
    </div>
  );
};

export default Seeds;
