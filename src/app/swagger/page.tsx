'use client';

import SwaggerUI from 'swagger-ui-react';
import 'swagger-ui-react/swagger-ui.css';

export default function SwaggerPage() {
  return (
    <main className="min-h-screen bg-slate-50 px-4 py-8">
      <div className="mx-auto max-w-7xl rounded-2xl bg-white shadow-lg">
        <SwaggerUI url="/api/openapi" />
      </div>
    </main>
  );
}
