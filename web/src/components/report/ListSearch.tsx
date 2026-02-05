import React, { useState, useEffect, useRef } from 'react';
import { levenshteinDistance } from '../../utils/levenshtein';

interface ListSearchProps<T> {
  items: T[];
  getLabel: (item: T) => string;
  onSelect: (item: T) => void;
  filterFn?: (query: string, items: T[], getLabel: (item: T) => string) => T[];
  placeholder?: string;
}

export function ListSearch<T>({ items, getLabel, onSelect, filterFn, placeholder }: ListSearchProps<T>) {
  const [query, setQuery] = useState('');
  const [selectedIndex, setSelectedIndex] = useState(0);
  const inputRef = useRef<HTMLInputElement>(null);

  const filtered = filterFn 
    ? filterFn(query, items, getLabel)
    : query
    ? items
        .map(item => {
          const label = getLabel(item).toLowerCase();
          const q = query.toLowerCase();
          const startsWithMatch = label.startsWith(q);
          const includesMatch = label.includes(q);
          
          // Check aliases if available
          const aliases = (item as any).aliases || [];
          const aliasStartsWith = aliases.some((a: string) => a.toLowerCase().startsWith(q));
          const aliasIncludes = aliases.some((a: string) => a.toLowerCase().includes(q));
          
          const distance = levenshteinDistance(q, label);
          return { item, startsWithMatch, includesMatch, aliasStartsWith, aliasIncludes, distance };
        })
        .sort((a, b) => {
          if (a.aliasStartsWith && !b.aliasStartsWith) return -1;
          if (!a.aliasStartsWith && b.aliasStartsWith) return 1;
          if (a.startsWithMatch && !b.startsWithMatch) return -1;
          if (!a.startsWithMatch && b.startsWithMatch) return 1;
          if (a.aliasIncludes && !b.aliasIncludes) return -1;
          if (!a.aliasIncludes && b.aliasIncludes) return 1;
          if (a.includesMatch && !b.includesMatch) return -1;
          if (!a.includesMatch && b.includesMatch) return 1;
          return a.distance - b.distance;
        })
        .map(({ item }) => item)
    : items;

  useEffect(() => {
    setSelectedIndex(0);
  }, [query]);

  useEffect(() => {
    inputRef.current?.focus();
  }, []);

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'ArrowDown') {
      e.preventDefault();
      e.stopPropagation();
      setSelectedIndex(i => Math.min(i + 1, filtered.length - 1));
    } else if (e.key === 'ArrowUp') {
      e.preventDefault();
      e.stopPropagation();
      setSelectedIndex(i => Math.max(i - 1, 0));
    } else if (e.key === 'Enter' && filtered.length > 0) {
      e.preventDefault();
      e.stopPropagation();
      onSelect(filtered[selectedIndex]);
    } else if (e.key === 'Escape') {
      e.preventDefault();
      e.stopPropagation();
    }
  };

  return (
    <div className="list-search">
      <input
        ref={inputRef}
        type="text"
        value={query}
        onChange={e => setQuery(e.target.value)}
        onKeyDown={handleKeyDown}
        placeholder={placeholder || 'Search...'}
      />
      <div className="list-search-results">
        {filtered.map((item, i) => (
          <div
            key={i}
            className={`list-search-item ${i === selectedIndex ? 'selected' : ''}`}
            onClick={() => onSelect(item)}
          >
            {getLabel(item)}
          </div>
        ))}
      </div>
    </div>
  );
}
