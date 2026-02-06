import React, { useState, useEffect } from 'react';
import { useTournament } from '../utils/tournament';
import API_URL, { getAuthHeaders } from '../config';
import { Button } from './common/Button';
import { DataTable } from './common/DataTable';
import { LoadingSpinner, ErrorMessage } from './common/LoadingError';
import { SeedResult } from '../types';

const Seeds: React.FC = () => {
  const [seeds, setSeeds] = useState<SeedResult[]>([]);
  const [tournament] = useTournament();
  const [redemption, setRedemption] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const runSeeding = async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await fetch(`${API_URL}/api/seed`, {
        method: 'POST',
        headers: getAuthHeaders(),
        body: JSON.stringify({ tournament, redemption }),
      });
      if (!response.ok) throw new Error('Failed to run seeding');
      const data = await response.json();
      setSeeds(data.seeds || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred');
    } finally {
      setLoading(false);
    }
  };

  const publishSeeding = async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await fetch(`${API_URL}/api/seed/publish`, {
        method: 'POST',
        headers: getAuthHeaders(),
        body: JSON.stringify({ tournament, redemption, players: seeds }),
      });
      if (!response.ok) throw new Error('Failed to publish seeds');
      alert('Seeds published successfully!');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    setSeeds([]);
  }, [tournament]);

  if (loading) return <LoadingSpinner />;

  return (
    <div className="page-container">
      <h1>Seeding</h1>
      {error && <ErrorMessage message={error} />}
      <div style={{ marginBottom: '20px', display: 'flex', gap: '20px', alignItems: 'center' }}>
        <label>
          <input
            type="checkbox"
            checked={redemption}
            onChange={e => setRedemption(e.target.checked)}
          />
          {' '}Redemption Bracket
        </label>
        <Button onClick={runSeeding}>Generate Seeds</Button>
      </div>

      {seeds.length > 0 && (
        <>
          <h2>Seeds ({seeds.length})</h2>
          <DataTable
            data={seeds}
            columns={[
              { key: 'seed', label: 'Seed' },
              { key: 'name', label: 'Name' },
              { 
                key: 'rating', 
                label: 'Rating',
                render: (s) => s.rating.toFixed(2)
              }
            ]}
            keyExtractor={(s) => s.id}
          />
          <Button variant="success" onClick={publishSeeding} style={{ marginTop: '20px' }}>
            Publish Seeds
          </Button>
        </>
      )}
    </div>
  );
};

export default Seeds;
